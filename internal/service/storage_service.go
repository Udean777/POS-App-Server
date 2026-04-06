package service

import (
	"bytes"
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/sajudin/pos-app-server/internal/domain"
)

type s3StorageService struct {
	client     *s3.Client
	bucketName string
	publicURL  string
}

func NewS3StorageService(ctx context.Context) (domain.StorageService, error) {
	accountID := os.Getenv("R2_ACCOUNT_ID")
	accessKey := os.Getenv("R2_ACCESS_KEY_ID")
	secretKey := os.Getenv("R2_SECRET_ACCESS_KEY")
	bucketName := os.Getenv("R2_BUCKET_NAME")
	publicURL := os.Getenv("R2_PUBLIC_DOMAIN")

	if accountID == "" || accessKey == "" || secretKey == "" || bucketName == "" {
		return nil, fmt.Errorf("R2 credentials missing")
	}

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
		config.WithRegion("auto"),
	)
	if err != nil {
		return nil, err
	}

	// Cloudflare R2 requires a custom endpoint
	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(fmt.Sprintf("https://%s.r2.cloudflarestorage.com", accountID))
	})

	return &s3StorageService{
		client:     client,
		bucketName: bucketName,
		publicURL:  publicURL,
	}, nil
}

func (s *s3StorageService) Upload(ctx context.Context, file []byte, fileName string, contentType string) (string, error) {
	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucketName),
		Key:         aws.String(fileName),
		Body:        bytes.NewReader(file),
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return "", err
	}

	// Ensure publicURL doesn't end with slash
	return fmt.Sprintf("%s/%s", s.publicURL, fileName), nil
}

type localStorageService struct {
	publicURL string
}

func NewLocalStorageService(publicURL string) domain.StorageService {
	return &localStorageService{publicURL: publicURL}
}

func (s *localStorageService) Upload(ctx context.Context, file []byte, fileName string, contentType string) (string, error) {
	if err := os.MkdirAll("uploads", 0755); err != nil {
		return "", err
	}
	if err := os.WriteFile("uploads/"+fileName, file, 0644); err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/uploads/%s", s.publicURL, fileName), nil
}
