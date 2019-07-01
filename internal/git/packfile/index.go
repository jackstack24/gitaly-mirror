package packfile

import (
	"bufio"
	"crypto/sha1"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"math"
	"os"
	"regexp"
	"sort"
	"strconv"

	"gitlab.com/gitlab-org/gitaly/internal/git/gitio"
)

const sumSize = sha1.Size

var (
	idxFileRegex = regexp.MustCompile(`\A(.*/pack-)([0-9a-f]{40})\.idx\z`)
)

type Index struct {
	ID       string
	packBase string
	Objects  []*Object
	fanOut   [256]int
	*Bitmap
}

func ReadIndex(idxPath string) (*Index, error) {
	reMatches := idxFileRegex.FindStringSubmatch(idxPath)
	if len(reMatches) == 0 {
		return nil, fmt.Errorf("invalid idx filename: %q", idxPath)
	}

	idx := &Index{
		packBase: reMatches[1] + reMatches[2],
		ID:       reMatches[2],
	}

	f, err := os.Open(idx.packBase + ".idx")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r := bufio.NewReader(gitio.NewHashfileReader(f))

	const sig = "\377tOc\x00\x00\x00\x02"
	actualSig, err := readN(r, len(sig))
	if s := string(actualSig); s != sig {
		return nil, fmt.Errorf("unexpected idx signature %q", s)
	}

	count, err := idx.nPackObjects()
	if err != nil {
		return nil, err
	}

	// TODO use a data structure other than a Go slice to hold the index
	// entries? Go slices use int as their index type, and int may not be
	// able to hold MaxUint32.
	if count > math.MaxInt32 {
		return nil, fmt.Errorf("too many objects in to fit in Go slice: %d", count)
	}
	idx.Objects = make([]*Object, count)

	for i := range idx.fanOut {
		n, err := readUint32(r)
		if err != nil {
			return nil, err
		}

		idx.fanOut[i] = int(n) // cast is safe because we know n<=len(idx.Objects)
	}

	buf := make([]byte, sumSize)
	for i := 0; i < len(idx.Objects); i++ {
		if _, err := io.ReadFull(r, buf); err != nil {
			return nil, err
		}
		idx.Objects[i] = &Object{OID: hex.EncodeToString(buf)}
	}

	// Discard CRC32 values (one for each object)
	for i := 0; i < len(idx.Objects); i++ {
		if _, err := r.Discard(4); err != nil {
			return nil, err
		}
	}

	// Read 4-byte offsets
	has8ByteOffsets := false
	for i := 0; i < len(idx.Objects); i++ {
		offset, err := readUint32(r)
		if err != nil {
			return nil, err
		}

		const mask = 1 << 31
		if offset&mask == mask {
			has8ByteOffsets = true
			continue
		}

		idx.Objects[i].Offset = uint64(offset)
	}

	if has8ByteOffsets {
		for i := 0; i < len(idx.Objects); i++ {
			offset, err := readUint64(r)
			if err != nil {
				return nil, err
			}

			// TODO Not clear if all 8-byte offsets are populated, or only those that
			// don't fit into 4 bytes.
			if offset > 0 {
				idx.Objects[i].Offset = offset
			}
		}
	}

	idxPackID, err := readN(r, sumSize)
	if err != nil {
		return nil, err
	}

	if s := hex.EncodeToString(idxPackID); s != idx.ID {
		return nil, fmt.Errorf("unexpected pack ID in idx: %s", s)
	}

	if _, err := r.Peek(1); err != io.EOF {
		if err == nil {
			err = fmt.Errorf("unexpected trailing data, expected EOF")
		}
		return nil, err
	}

	return idx, nil
}

func (idx *Index) GetObject(oid string) (*Object, bool) {
	if len(oid) < 2 {
		return nil, false
	}

	radix64, err := strconv.ParseInt(oid[:2], 16, 0)
	if err != nil {
		return nil, false
	}

	radix := int(radix64)
	last := idx.fanOut[radix]
	first := 0
	if radix > 0 {
		first = idx.fanOut[radix-1]
	}

	objRange := idx.Objects[first:last]
	objIdx := sort.Search(len(objRange), func(i int) bool {
		return objRange[i].OID >= oid
	})
	if objIdx == len(objRange) {
		return nil, false
	}
	obj := objRange[objIdx]

	if obj.OID != oid {
		return nil, false
	}

	return obj, true
}

func (idx *Index) nPackObjects() (uint32, error) {
	f, err := idx.openPack()
	if err != nil {
		return 0, err
	}
	defer f.Close()

	const headerLen = 12
	header, err := readN(f, headerLen)
	if err != nil {
		return 0, err
	}

	const sig = "PACK\x00\x00\x00\x02"
	if s := string(header[:len(sig)]); s != sig {
		return 0, fmt.Errorf("unexpected pack signature %q", s)
	}
	header = header[len(sig):]

	return binary.BigEndian.Uint32(header), nil
}

func (idx *Index) openPack() (f *os.File, err error) {
	packPath := idx.packBase + ".pack"
	f, err = os.Open(packPath)
	if err != nil {
		return nil, err
	}

	defer func(f *os.File) {
		if err != nil {
			f.Close()
		}
	}(f) // Bind f early so that we can do "return nil, err".

	if _, err := f.Seek(-sumSize, io.SeekEnd); err != nil {
		return nil, err
	}

	sum, err := readN(f, sumSize)
	if err != nil {
		return nil, err
	}

	if s := hex.EncodeToString(sum); s != idx.ID {
		return nil, fmt.Errorf("unexpected trailing checksum in .pack: %s", s)
	}

	if _, err := f.Seek(0, io.SeekStart); err != nil {
		return nil, err
	}

	return f, nil
}

func readUint32(r io.Reader) (uint32, error) {
	buf, err := readN(r, 4)
	if err != nil {
		return 0, err
	}

	return binary.BigEndian.Uint32(buf), nil
}

func readUint64(r io.Reader) (uint64, error) {
	buf, err := readN(r, 8)
	if err != nil {
		return 0, err
	}

	return binary.BigEndian.Uint64(buf), nil
}

func readN(r io.Reader, n int) ([]byte, error) {
	buf := make([]byte, n)
	if _, err := io.ReadFull(r, buf); err != nil {
		return nil, err
	}

	return buf, nil
}