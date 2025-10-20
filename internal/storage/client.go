package storage

import "os"

var (
	endpoint = os.Getenv("STORAGE_ENDPOINT")
	accessKeyId = os.Getenv("STORAGE_ACCESSKEY")
	secretAccessKey = os.Getenv("STORAGE_SECRET_ACCESSKEY")
	useSSL = true
)

func New() {

}