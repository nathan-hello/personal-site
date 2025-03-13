package render

import "strings"

func EscapeHtml(s string) string {

        s = strings.ReplaceAll(s, "&", "&amp;")
        s= strings.ReplaceAll(s, "<", "&lt;")
        s= strings.ReplaceAll(s, ">", "&gt;")
        s= strings.ReplaceAll(s, "\"", "&quot;")
        s= strings.ReplaceAll(s, "'", "&#039;")

        return s
}

