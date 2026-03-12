package components

import "unicode"

func sentenceizeString(s string) string {
	if s == "" {
		return ""
	}
	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	runes = append(runes, '.')
	return string(runes)
}

