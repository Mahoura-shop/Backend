package mocks

import (
	"mime/multipart"
	"time"

	"github.com/Mahoura-shop/Backend/internal/domain/enum"
	"github.com/stretchr/testify/mock"
)

type S3StorageMock struct {
	mock.Mock
}

func NewS3StorageMock() *S3StorageMock {
	return &S3StorageMock{}
}

func (m *S3StorageMock) DeleteObject(bucketType enum.BucketType, key string) error {
	args := m.Called(bucketType, key)
	return args.Error(0)
}

func (m *S3StorageMock) GetPresignedURL(bucketType enum.BucketType, objectKey string, expiration time.Duration) (string, error) {
	args := m.Called(bucketType, objectKey, expiration)
	return args.String(0), args.Error(1)
}

func (m *S3StorageMock) UploadObject(bucketType enum.BucketType, key string, file *multipart.FileHeader) error {
	args := m.Called(bucketType, key, file)
	return args.Error(0)
}
