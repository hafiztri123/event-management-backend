package storage

import (
	"context"
	"mime/multipart"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type StorageService interface {
	UploadFile(ctx context.Context, file multipart.File, filename string) (string, error)
	DeleteFile(ctx context.Context, publicID string) error
}

type cloudinaryService struct {
	cld *cloudinary.Cloudinary
}

func NewCloudinaryService(cloudName, apiKey, apiSecret string) (StorageService, error) {
	cld, err := cloudinary.NewFromParams(cloudName, apiKey, apiSecret)
	if err != nil {
		return nil, err
	}

	return &cloudinaryService{
		cld: cld,
	}, nil
}

func (s *cloudinaryService) UploadFile(ctx context.Context, file multipart.File, filename string) (string, error) {
	uploadResult, err := s.cld.Upload.Upload(ctx, file, uploader.UploadParams{
		PublicID: filename,
		ResourceType: "auto",
	})

	if err != nil {
		return "", err
	}

	return uploadResult.SecureURL, nil
}


func (s *cloudinaryService) DeleteFile(ctx context.Context, publicID string) error{
	_, err := s.cld.Upload.Destroy(ctx, uploader.DestroyParams{
		PublicID: publicID,
	})
	return err
}

