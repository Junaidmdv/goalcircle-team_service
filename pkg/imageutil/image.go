package imageutil

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"net/http"

	"github.com/Junaidmdv/goalcircle-team_service/pkg/apperror"
	"github.com/chai2010/webp"
)

type ImageType string

const (
	TeamLogo    ImageType = "TEAM_LOGO"
	PlayerImage ImageType = "PLAYER_IMAGE"
	StaffImage  ImageType = "STAFF_IMAGE"
)

type ImageConfig struct {
	MaxWidth     int
	MaxHeight    int
	MinWidth     int
	MinHeight    int
	MaxSize      int64
	AllowedTypes map[string]struct{}
}

var ImageConfigs = map[ImageType]ImageConfig{
	TeamLogo: {
		MinWidth:  256,
		MinHeight: 256,
		MaxWidth:  1024,
		MaxHeight: 1024,
		AllowedTypes: map[string]struct{}{
			"image/png":  {},
			"image/jpeg": {},
			"image/webp": {},
		},
	},
	PlayerImage: {
		MinWidth:  256,
		MinHeight: 256,
		MaxWidth:  1024,
		MaxHeight: 1024,
		AllowedTypes: map[string]struct{}{
			"image/png":  {},
			"image/jpeg": {},
			"image/webp": {},
		},
	}, StaffImage: {
		MinWidth:  256,
		MinHeight: 256,
		MaxWidth:  1024,
		MaxHeight: 1024,
		AllowedTypes: map[string]struct{}{
			"image/png":  {},
			"image/jpeg": {},
			"image/webp": {},
		},
	},
}

func ValidateImage(imageData []byte, imagetype ImageType) error {
	contentType := http.DetectContentType(imageData)

	m, ok := ImageConfigs[imagetype]
	if !ok {
		return  apperror.NewInternalError(apperror.InternalErrorMsg, errors.New("unknown image type"))

	}

	_, exist := m.AllowedTypes[contentType]

	if !exist {
		return  apperror.NewInvalidArgumentError("unsupported image content type")
	}

	config, _, err := image.DecodeConfig(bytes.NewReader(imageData))
	if err != nil {
		return  apperror.NewInternalError("failed decode image", err)
	}

	if config.Height > m.MaxHeight || config.Width > m.MaxWidth {
		return  apperror.NewInvalidArgumentError("logo dimensions must be at least 512x512")
	}

	if config.Height < m.MinHeight || config.Width < m.MinWidth {
		return  apperror.NewInvalidArgumentError("logo dimensions must be at least 256x256")
	}

	return  nil
}

func ConvertImageIntoWebpbFormate(imageData []byte) ([]byte, error) {
	img, _, err := image.Decode(bytes.NewReader(imageData))
	if err != nil {
		return nil, apperror.NewInternalError(apperror.InternalErrorMsg, fmt.Errorf("failed decode image to webp formate:%v", err))
	}

	var buf bytes.Buffer
	if err := webp.Encode(&buf, img, &webp.Options{Quality: 80}); err != nil {
		return nil, apperror.NewInternalError(apperror.InternalErrorMsg, fmt.Errorf("failed encode image to webp formate:%v", err))
	}

	return buf.Bytes(), nil

}
