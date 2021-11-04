package aws

import (
	"context"
	"io"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

var (
	bucketName = os.Getenv("AWS_S3_BUCKET")
	region = os.Getenv("AWS_REGION")
	accessKeyId = os.Getenv("AWS_ACCESS_KEY_ID")
	secretAccessKey = os.Getenv("AWS_SECRET_ACCESS_KEY")
)

type Bucket struct {
	name string
	uploader *s3manager.Uploader
}

func (b Bucket) Put(ctx context.Context, r io.Reader, key string) error {
	_, err := b.uploader.UploadWithContext(ctx, &s3manager.UploadInput{
		Bucket: aws.String(b.name),
		Key: aws.String(key),
		Body: r,
	})
	return err
}

func NewBucket() (Bucket, error) {
	if bucketName == "" {
		return Bucket{}, fmt.Errorf("new bucket: missing AWS_S3_BUCKET from env")
	}
	creds := credentials.NewStaticCredentials(accessKeyId, secretAccessKey, "")
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
		Credentials: creds,
	})
	if err != nil {
		return Bucket{}, err
	}

	return Bucket{
		name: bucketName,
		uploader: s3manager.NewUploader(sess),
	}, nil
}
