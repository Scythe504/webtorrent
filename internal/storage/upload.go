package storage

import (
	"context"
	"github.com/anacrolix/torrent"
	"github.com/minio/minio-go/v7"
)

func (s *service) WriteStream(ctx context.Context, id, objectName string, reader torrent.Reader) (string, error) {
	bucketName := "webtorrent-test"
	info, err := s.storage.PutObject(ctx, bucketName, "test.mkv", reader, -1, minio.PutObjectOptions{})

	if err != nil {
		return "", err
	}

	return info.Key, nil
}
