package storage

import (
	"fmt"
	"github.com/scythe504/webtorrent/internal/tor"
	"io"
	"os"
	"path/filepath"
)

// SaveForLater saves a torrent's video file using metadata to determine filename & extension.
func (s *service) SaveForLater(videoId string, reader io.Reader, meta tor.FileMetadata) (string, error) {
	if !meta.IsVideo {
		return "", fmt.Errorf("file %s is not a recognized video type", meta.Name)
	}

	fileName := fmt.Sprintf(meta.Name)
	targetPath := filepath.Join(s.dataDir, fileName)

	file, err := os.Create(targetPath)
	if err != nil {
		return "", fmt.Errorf("failed to create file %s: %w", targetPath, err)
	}
	defer file.Close()

	if _, err := io.Copy(file, reader); err != nil {
		return "", fmt.Errorf("failed to write file data for %s: %w", videoId, err)
	}

	return targetPath, nil
}
