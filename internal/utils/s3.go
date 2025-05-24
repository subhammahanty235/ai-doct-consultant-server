package utils

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
)

type S3Client struct {
	session *session.Session
	svc     *s3.S3
	bucket  string
}

func NewS3Client(region, accessKey, secretKey, bucket string) (*S3Client, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
		Credentials: credentials.NewStaticCredentials(
			accessKey,
			secretKey,
			"",
		),
	})
	if err != nil {
		return nil, err
	}

	return &S3Client{
		session: sess,
		svc:     s3.New(sess),
		bucket:  bucket,
	}, nil
}

func (s *S3Client) UploadFile(file multipart.File, fileHeader *multipart.FileHeader) (string, error) {
	// Generate unique filename
	ext := filepath.Ext(fileHeader.Filename)
	filename := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	key := fmt.Sprintf("chat-images/%s", filename)

	// Read file content
	buffer := make([]byte, fileHeader.Size)
	file.Read(buffer)

	// Upload to S3
	_, err := s.svc.PutObject(&s3.PutObjectInput{
		Bucket:        aws.String(s.bucket),
		Key:           aws.String(key),
		Body:          bytes.NewReader(buffer),
		ContentLength: aws.Int64(fileHeader.Size),
		ContentType:   aws.String(fileHeader.Header.Get("Content-Type")),
		ACL:           aws.String("public-read"),
	})

	if err != nil {
		return "", err
	}

	// Return public URL
	url := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", s.bucket, key)
	return url, nil
}

func (s *S3Client) GeneratePresignedURL(key string) (string, error) {
	req, _ := s.svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})

	url, err := req.Presign(15 * time.Minute)
	if err != nil {
		return "", err
	}

	return url, nil
}
