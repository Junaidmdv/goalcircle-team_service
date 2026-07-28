package team 



import (
	"strings"

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