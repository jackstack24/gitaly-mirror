package praefect

import (
	"bytes"
	"time"

	"gitlab.com/gitlab-org/gitaly/proto/go/gitalypb"
)

const base62Chars string = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// encodeReverseBase62 encodes num into its Base62 reversed representation.
// The most significant value is at the end of the string.
//
// Appending is faster than prepending and this is enough for the purpose of a random ID
func encodeReverseBase62(num int64) string {
	if num == 0 {
		return "0"
	}

	encoded := bytes.Buffer{}
	for q := num; q > 0; q /= 62 {
		encoded.Write([]byte{base62Chars[q%62]})
	}

	return encoded.String()
}

// generates a correlation ID for each repo and time combination
func generatePseudorandomCorrelationID(repo *gitalypb.Repository) string {
	return repo.GetRelativePath() + ":" + encodeReverseBase62(time.Now().UnixNano())
}
