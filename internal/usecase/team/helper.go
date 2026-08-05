package team

import (
	"strings"

	"github.com/Junaidmdv/goalcircle-team_service/pkg/apperror"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
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
