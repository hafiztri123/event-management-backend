package storage

import (
	"context"
	"mime/multipart"
	"strings"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/hafiztri123/src/internal/pkg/config"
)

type StorageService interface {
	UploadFile(ctx context.Context, file multipart.File, filename string) (string, error)
	DeleteFile(ctx context.Context, publicID string) error
}

type cloudinaryService struct {
	cld *cloudinary.Cloudinary
}

func NewCloudinaryService(cfg *config.Config) (StorageService, error) {
	cld, err := cloudinary.NewFromParams(cfg.CloudinaryConfig.CloudName, cfg.CloudinaryConfig.ApiKey, cfg.CloudinaryConfig.ApiSecret)
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


// Example URL: https://res.cloudinary.com/your-cloud-name/image/upload/v1234567890/folder/filename.jpg

func ExtractPublicID(url string) string {
	if url == ""{
		return ""
	}

	parts := strings.Split(url, "/")

	uploadIndex := -1
	for i, part := range parts {
		if part == "upload" {
			uploadIndex = i
			break
		}
	}

	if uploadIndex == -1 || uploadIndex+2 >= len(parts) {
		return ""
	}

	publicID := strings.Join(parts[uploadIndex+2:], "/")
	
	if dotIndex := strings.LastIndex(publicID, "."); dotIndex != -1 {
		publicID = publicID[:dotIndex]
	}

	return publicID
	
}

