package minio

import (
	"context"
	"os"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func NewMinioClient(endpoint, id, secret string) (*minio.Client, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(id, secret, ""),
		Secure: false,
	})
	if err != nil {
		return nil, err
	}
	return client, nil
}

func FileUpload(client *minio.Client, bucket, name, filePath string) error {
	ctx := context.Background()
	exists, err := client.BucketExists(ctx, bucket)
	if err != nil {
		return err
	}
	if !exists {
		if err = client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{}); err != nil {
			return err
		}
	}
	obj, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer obj.Close()
	stat, err := obj.Stat()
	if err != nil {
		return err
	}
	_, err = client.PutObject(ctx, bucket, name, obj, stat.Size(), minio.PutObjectOptions{ContentType: "application/octet-stream"})
	if err != nil {
		return err
	}
	return nil
}

func FileDownload(client *minio.Client, bucket, name, filePath string) error {
	ctx := context.Background()
	err := client.FGetObject(ctx, bucket, name, filePath, minio.GetObjectOptions{})
	if err != nil {
		return err
	}
	return nil
}

func FilePreview(client *minio.Client, bucket, name string, expires time.Duration) (string, error) {
	ctx := context.Background()
	u, err := client.PresignedGetObject(ctx, bucket, name, expires, nil)
	if err != nil {
		return "", err
	}
	return u.String(), nil
}
