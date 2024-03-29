package storage

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"gitlab.com/gitlab-org/gitaly/internal/config"
	"gitlab.com/gitlab-org/gitaly/internal/safe"
)

const (
	// metadataFilename is the filename for a file we write on the gitaly server containing metadata about
	// the filesystem
	metadataFilename = ".gitaly-metadata"
)

// Metadata contains metadata about the filesystem
type Metadata struct {
	GitalyFilesystemID string `json:"gitaly_filesystem_id"`
}

// WriteMetadataFile marshals and writes a metadata file
func WriteMetadataFile(storage config.Storage) error {
	path := filepath.Join(storage.Path, metadataFilename)

	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return err
	}

	fw, err := safe.CreateFileWriter(path)
	if err != nil {
		return err
	}
	defer fw.Close()

	if err = json.NewEncoder(fw).Encode(&Metadata{
		GitalyFilesystemID: uuid.New().String(),
	}); err != nil {
		return err
	}

	return fw.Commit()
}

// ReadMetadataFile reads and decodes the json metadata file
func ReadMetadataFile(storage config.Storage) (Metadata, error) {
	path := filepath.Join(storage.Path, metadataFilename)

	var metadata Metadata

	metadataFile, err := os.Open(path)
	if err != nil {
		return metadata, err
	}
	defer metadataFile.Close()

	if err = json.NewDecoder(metadataFile).Decode(&metadata); err != nil {
		return metadata, err
	}

	return metadata, nil
}
