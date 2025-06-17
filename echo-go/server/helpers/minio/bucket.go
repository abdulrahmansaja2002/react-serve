package minio

import (
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"strings"

	"github.com/minio/minio-go/v7"
)

type Config struct {
	ErrorOnExisting bool
	PublicAccess    bool
	Region          string
}

func checkBucketExists(ctx context.Context, mc *minio.Client, name string) bool {
	exists, err := mc.BucketExists(ctx, name)
	if err != nil {
		log.Printf("Error checking if bucket exists: %s\n", err)
		return false
	}
	return exists
}

func applyPublicAccessPolicy(ctx context.Context, mc *minio.Client, name string) {
	policy := `{
		"Version": "2012-10-17",
		"Statement": [
			{
				"Effect": "Allow",
				"Principal": {
					"AWS": ["*"]
				},
				"Action": [
					"s3:GetObject"
				],
				"Resource": [
					"arn:aws:s3:::` + name + `/*"
				]
			}
		]
	}`
	err := mc.SetBucketPolicy(ctx, name, policy)
	if err != nil {
		log.Printf("Error setting bucket policy: %s\n", err)
	}
}

func CreateBucket(ctx context.Context, mc *minio.Client, name string) {
	if checkBucketExists(ctx, mc, name) {
		log.Printf("Bucket %s already exists\n", name)
		return
	}
	err := mc.MakeBucket(ctx, name, minio.MakeBucketOptions{Region: "us-east-1"})
	if err != nil {
		log.Printf("Error creating bucket: %s\n", err)
		return
	}
	// applyPublicAccessPolicy(ctx, mc, name)
	log.Printf("Bucket %s created successfully\n", name)
}

func ToObjectPath(bucket string, object string) string {
	return fmt.Sprintf("%s/%s", bucket, object)
}

func extractPath(path string) (bucket string, object string) {
	trimPath := strings.TrimLeft(path, "/")
	splitedPath := strings.Split(trimPath, "/")

	if len(splitedPath) < 2 {
		return "", ""
	}
	bucket = splitedPath[0]
	object = strings.Join(splitedPath[1:], "/")
	return
}

func PutObject(ctx context.Context, mc *minio.Client, path string, fileIn *multipart.FileHeader) (err error) {
	bucket, name := extractPath(path)
	if bucket == "" || name == "" {
		return fmt.Errorf("invalid path: %s", path)
	}
	file, err := fileIn.Open()
	if err != nil {
		return
	}
	defer file.Close()
	info, err := mc.PutObject(ctx, bucket, name, file, fileIn.Size, minio.PutObjectOptions{ContentType: fileIn.Header.Get("Content-Type")})
	if err != nil {
		return
	}
	log.Println(info)
	return
}

func GetObject(ctx context.Context, mc *minio.Client, path string) (object *minio.Object, err error) {
	bucket, name := extractPath(path)
	if bucket == "" || name == "" {
		return nil, fmt.Errorf("invalid path: %s", path)
	}
	return mc.GetObject(ctx, bucket, name, minio.GetObjectOptions{})
}

func DeleteObject(ctx context.Context, mc *minio.Client, path string) (err error) {
	bucket, name := extractPath(path)
	if bucket == "" || name == "" {
		return fmt.Errorf("invalid path: %s", path)
	}
	err = mc.RemoveObject(ctx, bucket, name, minio.RemoveObjectOptions{})
	if err != nil {
		return
	}
	return
}

func GetObjectInfo(ctx context.Context, mc *minio.Client, path string) (info minio.ObjectInfo, err error) {
	bucket, name := extractPath(path)
	if bucket == "" || name == "" {
		return minio.ObjectInfo{}, fmt.Errorf("invalid path: %s", path)
	}
	return mc.StatObject(ctx, bucket, name, minio.StatObjectOptions{})
}

func GetListObjects(ctx context.Context, mc *minio.Client, bucket string) (info []minio.ObjectInfo, err error) {
	// Inisialisasi slice untuk menampung hasil
	var results []minio.ObjectInfo

	if bucket == "" {
		return nil, fmt.Errorf("invalid bucket: %s", bucket)
	}

	// Dapatkan channel object dari MinIO
	objectCh := mc.ListObjects(ctx, bucket, minio.ListObjectsOptions{Recursive: true})

	// Iterasi channel
	for obj := range objectCh {
		if obj.Err != nil {
			return nil, fmt.Errorf("error listing object %s: %w", obj.Key, obj.Err)
		}
		results = append(results, obj)
	}

	return results, nil
}

func GetObjects(ctx context.Context, mc *minio.Client, bucket string, key string) (object *minio.Object, err error) {
	if bucket == "" || key == "" {
		return nil, fmt.Errorf("invalid path: %s", key)
	}
	return mc.GetObject(ctx, bucket, key, minio.GetObjectOptions{})
}