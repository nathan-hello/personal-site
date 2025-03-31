package render

import (
	"fmt"
	"strings"
)

func EscapeHtml(s string) string {
	var buf strings.Builder
	for i := 0; i < len(s); {
		// >> case: link delimited by space (or non-digit)
		if i+1 < len(s) && s[i] == '>' && s[i+1] == '>' {
			j := i + 2
			// Collect digits only.
			for j < len(s) && s[j] >= '0' && s[j] <= '9' {
				j++
			}
			// If no digits were found, fall back to normal escaping.
			if j == i+2 {
				buf.WriteString(`<span class="text-lime-400">></span><span class="text-lime-400">></span>`)
				i += 2
				continue
			}
			token := s[i:j]
			// Delimiter is a space: if present, leave it out of the token.
			buf.WriteString(fmt.Sprintf(`<a href="#comment-%s">%s</a>`, token[2:], token))
			i = j
			continue
		}
		// > case: colored text delimited by newline (or the start of an HTML tag)
		if s[i] == '>' {
			j := i + 1
			for j < len(s) && s[j] != '\n' && s[j] != '<' {
				j++
			}
			token := s[i:j]
			buf.WriteString(fmt.Sprintf(`<span class="text-lime-400">%s</span>`, token))
			i = j
			continue
		}
		switch s[i] {
		case '<':
			buf.WriteString("<span><</span>")
		case '&':
			buf.WriteString("<span>&</span>")
		case '"':
			buf.WriteString("<span>\"</span>")
		case '\'':
			buf.WriteString("<span>'</span>")
		default:
			buf.WriteByte(s[i])
		}
		i++
	}
	return buf.String()
}
