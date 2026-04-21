package media

import (
	"context"
	"mime/multipart"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

// CloudinaryService wraps the Cloudinary SDK client
type CloudinaryService struct {
	client *cloudinary.Cloudinary
}

// NewCloudinaryService creates a new instance of our Cloudinary service
func NewCloudinaryService(secretURL string) (*CloudinaryService, error) {
	// Initialize the Cloudinary client using the URL from your .env file
	cld, err := cloudinary.NewFromURL(secretURL)
	if err != nil {
		return nil, err
	}
	
	return &CloudinaryService{client: cld}, nil
}

// UploadFile takes an uploaded file from the frontend and sends it to Cloudinary
func (s *CloudinaryService) UploadFile(ctx context.Context, file multipart.File) (string, error) {
	// Upload file to Cloudinary
	// We specify a folder name so all her birthday media is neatly organized in your Cloudinary dashboard!
	resp, err := s.client.Upload.Upload(ctx, file, uploader.UploadParams{
		Folder: "babe_25th_birthday", 
	})
	
	if err != nil {
		return "", err
	}
	
	// Return the secure HTTPS URL so we can save it in MongoDB
	return resp.SecureURL, nil
}




