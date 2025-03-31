package render

import "strings"

func EscapeHtml(s string) string {
    var buf strings.Builder
    for _, ch := range s {
        switch ch {
        case '&':
            buf.WriteString("<span>&</span>")
        case '<':
            buf.WriteString("<span><</span>")
        case '>':
            buf.WriteString("<span>></span>")
        case '"':
            buf.WriteString("<span>\"</span>")
        case '\'':
            buf.WriteString("<span>'</span>")
        default:
            buf.WriteRune(ch)
        }
    }
    return buf.String()
}
