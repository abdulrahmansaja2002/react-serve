package config

import (
	"context"

	"echo-react-serve/constants/bucket"
	minoHelper "echo-react-serve/helpers/minio"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var MinioClient *minio.Client

func InitMinioClient() {
	endpoint := Envs.Storage.Endpoint
	accessKeyID := Envs.Storage.AccessKey
	secretAccessKey := Envs.Storage.SecretKey
	useSSL := Envs.Storage.SslMode

	mc, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		panic(err)
	}
	MinioClient = mc
	insertAllBucket()
}

func insertAllBucket() {
	for _, b := range bucket.Buckets {
		minoHelper.CreateBucket(context.Background(), MinioClient, b)
	}
}
