package uploadtos3

import (
	dynamodbupload "awsapp/dynamodb"
	"awsapp/errors"
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func readContent(body io.Reader) string {
	bytes, err := io.ReadAll(body)
	errors.HandleError(err, "Reading content of S3 bucket record")
	return string(bytes)
}

func UploadFile(client *s3.Client, bucket, filePath, key string) error {

	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Cannot find the file specified. Error %v", err)
	}

	defer file.Close()

	_, err = client.PutObject(context.Background(), &s3.PutObjectInput{
		Bucket: &bucket,
		Key:    &key,
		Body:   file,
	})

	if err != nil {
		log.Fatalf("Cannot insert data into the S3 bucket. Error: %v", err)
	}

	return nil

}

func UploadToS3(filePath string) {

	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion("us-east-1"))
	if err != nil {
		log.Fatalf("Unable to AWS Default Config. Error %v", err)
	}

	s3Client := s3.NewFromConfig(cfg)
	bucketName := "nikhil1502"
	fileName := filePath
	key := "user1"

	err = UploadFile(s3Client, bucketName, fileName, key)
	fmt.Println("Upload Successful to S3 bucket")

	resp, err := s3Client.GetObject(context.Background(), &s3.GetObjectInput{
		Bucket: &bucketName,
		Key:    &key,
	})

	defer resp.Body.Close()

	text := readContent(resp.Body)
	dynamodbupload.StoreSummary(key, text)
}
