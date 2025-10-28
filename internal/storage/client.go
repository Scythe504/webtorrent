package storage

import (
	"os"

	"github.com/anacrolix/torrent"
)

type Service interface {
	SaveForLater(videoId string, reader torrent.Reader) error
	GetFilePath(videoId string) (string, error)
	FileExists(videoId string) bool
}

type service struct {
	dataDir string
}

func New() Service {
	dataDir := os.Getenv("DATA_DIR")
	if dataDir == "" {
		dataDir = "./download"
	}

	// Ensure directory exists
	os.MkdirAll(dataDir, 0755)

	return &service{
		dataDir: dataDir,
	}
}
