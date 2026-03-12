package render

import (
	"fmt"
	"strconv"
	"strings"
)

func EscapeHtml(s string) (string, []int64) {
	var buf strings.Builder
    inTag := false
    replies := []int64{}
	for i := 0; i < len(s); {
        if strings.HasPrefix(s[i:], "```") {
        	end := strings.Index(s[i+3:], "```")
        	if end != -1 {
        		end += i + 3
        		buf.WriteString(s[i : end+3])
        		i = end + 3
        		continue
        	}
        	buf.WriteString(s[i:])
        	break
        }
		if i+1 < len(s) && s[i] == '>' && s[i+1] == '>' {
			j := i + 2
			for j < len(s) && s[j] >= '0' && s[j] <= '9' {
				j++
			}
			if j == i+2 {
				buf.WriteString(`<span class="text-lime-400">></span><span class="text-lime-400">></span>`)
				i += 2
				continue
			}
			token := s[i:j]
			buf.WriteString(fmt.Sprintf(`<a href="#comment-%s">%s</a>`, token[2:], token))
            id, _ := strconv.Atoi(token[2:])
            replies = append(replies, int64(id))
			i = j
			continue
		}
		if s[i] == '<' {
			inTag = true
			buf.WriteString("<span><</span>")
			i++
			continue
		}

		if inTag {
			if s[i] == '>' {
				inTag = false
				buf.WriteByte('>')
				i++
				continue
			}
			buf.WriteByte(s[i])
			i++
			continue
		}

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
	return buf.String(), replies
}
