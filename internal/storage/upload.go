package storage

import (
	"fmt"
	"io"
	"log"
	"path/filepath"

	"github.com/anacrolix/torrent"
)

func (s *service) SaveForLater(videoId string, reader torrent.Reader) error {
	log.Printf("[SaveForLater] Starting full download for video: %s", videoId)

	// Read till EOF to force torrent to download entire file
	buf := make([]byte, 5*1024*1024) // 5MB buffer
	totalRead := int64(0)

	for {
		n, err := reader.Read(buf)
		totalRead += int64(n)

		if n > 0 {
			log.Printf("[SaveForLater] Downloaded: %d MB", totalRead/(1024*1024))
		}

		if err == io.EOF {
			log.Printf("[SaveForLater] Complete! Total: %d MB", totalRead/(1024*1024))
			break
		}

		if err != nil {
			return fmt.Errorf("download failed: %w", err)
		}
	}

	return nil
}

func (s *service) GetFilePath(videoId string) (string, error) {
	// Torrent files are saved by anacrolix with their original filename
	// You'll need to track the actual filename in your database
	// For now, pattern match to find the file

	pattern := filepath.Join(s.dataDir, "*")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return "", err
	}

	// This is simplified - you should store the actual filename in DB
	// and look it up by videoId
	if len(matches) > 0 {
		return matches[0], nil
	}

	return "", fmt.Errorf("file not found for video: %s", videoId)
}

func (s *service) FileExists(videoId string) bool {
	_, err := s.GetFilePath(videoId)
	return err == nil
}
