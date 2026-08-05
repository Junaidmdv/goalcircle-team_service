package team

import (
	"fmt"
	"strings"

	"github.com/Junaidmdv/goalcircle-team_service/pkg/apperror"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const (
	MinWidth  = 128
	MinHeight = 128

	MaxWidth  = 1024
	MaxHeight = 1024
)

func FormatTeamName(name string) string {
	// Trim and remove extra spaces
	name = strings.Join(strings.Fields(strings.TrimSpace(name)), " ")

	caser := cases.Title(language.English)
	name = caser.String(strings.ToLower(name))

	// Keep common abbreviations uppercase
	replacer := strings.NewReplacer(
		" Fc", " FC",
		" Sc", " SC",
		" Ac", " AC",
		" Cf", " CF",
		" Afc", " AFC",
		" Utd", " UTD", // or "United", depending on your preference
	)

	return replacer.Replace(name)
}

func ImageAllowedFormate(mimetype string) error {
	allowedImageFormate := map[string]struct{}{
		"image/jpg":  {},
		"image/webp": {},
		"image/png":  {},
	}

	if _, ok := allowedImageFormate[mimetype]; !ok {
		return apperror.NewInvalidArgumentError("unsupported image content type")
	}

	return nil
}

func ValidateImageDiamension(height, width int) error {

	if height < MinHeight || width < MinWidth {
		return apperror.NewInvalidArgumentError("logo dimensions must be at least 128x128")
	}

	if height > MaxHeight || width > MaxWidth {
		return apperror.NewInvalidArgumentError("logo dimensions must be at least 128x128")
	}
	return nil
}

func CreateObjectName(teamID string) string {
	return fmt.Sprintf("teams/%s/logo.webp", teamID)
}
