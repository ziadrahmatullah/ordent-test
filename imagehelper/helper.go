package imagehelper

import (
	"context"
	"io"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/ziadrahmatullah/ordent-test/config"
)

type ImageHelper interface {
	Upload(ctx context.Context, file io.Reader, folder string, key string) (string, error)
	Destroy(ctx context.Context, folder string, key string) error
}

func NewImageHelper() (ImageHelper, error) {
	cldConfig := config.NewCloudinaryConfig()
	cld, err := cloudinary.NewFromParams(cldConfig.Name, cldConfig.Key, cldConfig.Secret)
	return &cloudinaryService{cld: cld}, err
}
