package customs

import (
	"bytes"
	"strings"

	"github.com/a-h/templ"
	"github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
)

func code(c component) (templ.Component, error) {
	highlighted, err := CodeHighlighter(c.Attributes["lang"], c.Children)
	if err != nil {
		afterStart := strings.Split(c.Children, ">")[1]
		lts := strings.Split(afterStart, "<")
		beforeEnd := lts[len(lts)-2] // minus 2. 1 for length and 1 to have the second to last idx
		highlighted = beforeEnd
	}
	return templ.Raw(highlighted), nil
}

func CodeHighlighter(lang string, content string) (string, error) {
	lex := lexers.Get(lang)
	if lex == nil {
		lex = lexers.Analyse(lang)
	}
        if lex == nil {
                lex = lexers.Fallback
        }
	sty := styles.Get("gruvbox")
	frm := html.New(
		html.ClassPrefix("chroma-"),
		html.WithLineNumbers(true),
	)
	iter, err := lex.Tokenise(nil, content)
	if err != nil {
		return "", nil
	}
	var buf bytes.Buffer
	err = frm.Format(&buf, sty, iter)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
