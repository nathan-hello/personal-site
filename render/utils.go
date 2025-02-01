package render

import "strings"

func EscapeHtml(s string) string {

        strings.ReplaceAll(s, "&", "&amp;")
        strings.ReplaceAll(s, "<", "&lt;")
        strings.ReplaceAll(s, ">", "&gt;")
        strings.ReplaceAll(s, "\"", "&quot;")
        strings.ReplaceAll(s, "'", "&#039;")

        return s
}

