package storage

import (
	"context"
	"io"
	"log"

	"github.com/anacrolix/torrent"
	"github.com/minio/minio-go/v7"
)

func (s *service) WriteStream(ctx context.Context, id, objectName string, reader torrent.Reader) (string, error) {
	bucketName := "webtorrent-test"

	pipeReader, pipeWriter := io.Pipe()

	go func() {
		defer pipeWriter.Close()
		buf := make([]byte, 5*1024*1024) // 5MB chunks
		for {
			n, err := reader.Read(buf)
			if n > 0 {
				log.Printf("[WriteStream] Read %d bytes from torrent", n)
				pipeWriter.Write(buf[:n])
			}
			if err != nil {
				if err != io.EOF {
					log.Printf("[WriteStream] torrent read error: %v", err)
					pipeWriter.CloseWithError(err)
				}
				break
			}
		}
	}()

	info, err := s.storage.PutObject(ctx, bucketName, objectName, pipeReader, -1, minio.PutObjectOptions{})
	if err != nil {
		log.Printf("[WriteStream] MinIO upload error: %v", err)
		return "", err
	}

	return info.Key, nil
}

// Add this method to your storage service
func (s *service) StatObject(ctx context.Context, objectName string) (int64, error) {
	bucketName := "webtorrent-test"

	objInfo, err := s.storage.StatObject(ctx, bucketName, objectName, minio.StatObjectOptions{})
	if err != nil {
		return 0, err
	}

	return objInfo.Size, nil
}
