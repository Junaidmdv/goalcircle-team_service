package invitation

import (
	"crypto/rand"
	"math/big"
	"strings"
)

type CodeGenerater interface {
	GenerateShortName(string) string
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

func (cg *codeGenerater) GenerateShortName(teamName string) string {
	words := strings.Fields(strings.TrimSpace(teamName))

	switch len(words) {
	case 0:
		return ""
	case 1:
		name := strings.ToUpper(words[0])
		if len(name) >= 3 {
			return name[:3]
		}
		return name
	default:
		short := ""

		for _, word := range words {
			short += strings.ToUpper(string(word[0]))
		}

		if len(short) >= 3 {
			return short[:3]
		}

		name := strings.ToUpper(strings.ReplaceAll(teamName, " ", ""))
		for len(short) < 3 && len(name) > len(short) {
			short += string(name[len(short)])
		}

		return short
	}
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

	code := cg.GenerateShortName(role) + string(b)
	return code, nil
}
