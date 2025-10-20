package storage

import (
	"context"
	"log"
	"os"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var (
	endpoint        = os.Getenv("STORAGE_ENDPOINT")
	accessKeyId     = os.Getenv("STORAGE_ACCESSKEY")
	secretAccessKey = os.Getenv("STORAGE_SECRET_ACCESSKEY")
	useSSL          = true
)

type Service interface {
}

type service struct {
	storage *minio.Client
}

func New(ctx context.Context) Service {

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyId, secretAccessKey, ""),
		Secure: useSSL,
	})

	if err != nil {
		log.Fatal(err)
	}

	bucketName := "webtorrent-test"

	err = client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: ""})
	if err != nil {
		exists, errBucketExists := client.BucketExists(ctx, bucketName)
		if errBucketExists != nil && exists {
			log.Println("We already have a bucket with the name", bucketName)
		} else {
			log.Fatal(err)
		}
	} else {
		log.Printf("Bucket %q created successfully\n", bucketName)
	}

	st := service{storage: client}

	return st
}
