package invitation

import (
	"crypto/rand"
	"math/big"
	"regexp"
	"strings"
)

type CodeGenerater interface {
	GenerateShortName(string, int, bool) string
	GenerateCode(string) (string, error)
}

type codeGenerater struct {
	chars  string
	length int32
}

func NewCodeGenerater(length int32) CodeGenerater {
	chars := "ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

	return &codeGenerater{
		chars:  chars,
		length: length,
	}
}

var numberRegex = regexp.MustCompile(`^\d+$`)
var vowels = map[byte]bool{'a': true, 'e': true, 'i': true, 'o': true, 'u': true}

func firstLetter(w string) string {
	for i := 0; i < len(w); i++ {
		c := w[i]
		if c >= 'a' && c <= 'z' {
			return string(c)
		}
	}
	if len(w) > 0 {
		return string(w[0])
	}
	return ""
}

func (cg *codeGenerater) GenerateShortName(teamName string, shortWordMax int, uppercase bool) string {
	if shortWordMax <= 0 {
		shortWordMax = 2
	}

	rawWords := strings.Fields(strings.TrimSpace(teamName))
	words := make([]string, 0, len(rawWords))
	for _, w := range rawWords {
		if !numberRegex.MatchString(w) {
			words = append(words, w)
		}
	}

	if len(words) == 0 {
		return ""
	}

	var result string

	if len(words) == 1 {
		w := strings.ToLower(words[0])
		result = string(w[0])
		found := false
		for i := 1; i < len(w); i++ {
			if !vowels[w[i]] {
				result += string(w[i])
				found = true
				break
			}
		}
		// fallback if no consonant found (e.g. "aeiou" edge case)
		if !found && len(w) > 1 {
			result += string(w[1])
		}
	} else {
		// Multiple words -> keep short words whole, else take first letter
		var sb strings.Builder
		for _, w := range words {
			lw := strings.ToLower(w)
			if len(lw) <= shortWordMax {
				sb.WriteString(lw)
			} else {
				sb.WriteString(firstLetter(lw))
			}
		}
		result = sb.String()
	}

	if uppercase {
		return strings.ToUpper(result)
	}
	return result

}

func (cg *codeGenerater) GenerateCode(role string) (string, error) {

	b := make([]byte, cg.length)

	for i := range b {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(cg.chars))))
		if err != nil {
			return "", nil
		}

		b[i] = cg.chars[n.Int64()]
	}

	code := cg.GenerateShortName(role, 2, true) + string(b)
	return code, nil
}
