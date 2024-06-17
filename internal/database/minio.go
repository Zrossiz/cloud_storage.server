package db

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func InitMinio() *minio.Client {
	endpoint := os.Getenv("MINIO_ENDPOINT")
	accessKeyID := os.Getenv("MINIO_ACCESS_KEY")
	secretAccessKey := os.Getenv("MINIO_SECRET_KEY")
	useSSL := false

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalf("Failed to create MinIO client: %v", err)
	}

	bucketName := os.Getenv("BUCKET_NAME")
	location := "eu-central-1"

	err = createBucket(minioClient, bucketName, location)
	if err != nil {
		log.Fatalf("Failed to create bucket: %v", err)
	}

	return minioClient
}

func createBucket(client *minio.Client, bucketName, location string) error {
	ctx := context.Background()
	err := client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: location})
	if err != nil {
		exists, errBucketExists := client.BucketExists(ctx, bucketName)
		if errBucketExists == nil && exists {
			fmt.Printf("Bucket %s already exists\n", bucketName)
			return nil
		}
		return fmt.Errorf("failed to create bucket %v: %v", bucketName, err)
	}

	fmt.Printf("Successfully created bucket %s\n", bucketName)
	return nil
}
