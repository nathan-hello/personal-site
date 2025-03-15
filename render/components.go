package render

import (
	"bytes"
	"context"
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/a-h/templ"
	"github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
	"github.com/nathan-hello/personal-site/components"
	"github.com/nathan-hello/personal-site/utils"
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
	iter, err := lex.Tokenise(nil, content)
	if err != nil {
		return "", nil
	}

	// chroma-classes are written in ./public/css/chroma.css
	sty := styles.Get("gruvbox")
	frm := html.New(
		html.WrapLongLines(false),
		html.WithClasses(true),
		html.ClassPrefix("chroma-"),
		html.WithLineNumbers(true),
	)

	var buf bytes.Buffer
	err = frm.Format(&buf, sty, iter)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

// <ManyPostMini amount="all | <int>" author="nathan | natalie | all" sort="ascending | descending"/>
func manyPostMini(c component) (templ.Component, error) {

	amount, ok := c.Attributes["amount"]
	if !ok {
		return nil, fmt.Errorf("component did not give required attribute \"amount\": %#v", c)
	}
	author, ok := c.Attributes["author"]
	if !ok {
		return nil, fmt.Errorf("component did not give required attribute \"author\": %#v", c)
	}
	sort, ok := c.Attributes["sort"]
	if !ok {
		sort = "descending"
	}

	blogs, err := Blogs("./public/content/blog", "", false)
	if err != nil {
		return nil, err
	}

	amt, err := strconv.Atoi(amount)
	if err != nil {
		if amount == "all" {
			amt = 9999999
		} else {
			return nil, fmt.Errorf("amount has bad value: %#v", c)
		}
	}

	if sort == "descending" {
		slices.SortFunc(blogs, func(a, b utils.Blog) int {
			return b.Frnt.Date.Compare(a.Frnt.Date)
		})

	}
	if sort == "ascending" {
		slices.SortFunc(blogs, func(a, b utils.Blog) int {
			return a.Frnt.Date.Compare(b.Frnt.Date)
		})
	}

	if author != "all" {
		tmp := []utils.Blog{}
		for _, v := range blogs {
			if v.Frnt.Author == author {
				tmp = append(tmp, v)
			}
		}
		blogs = tmp
	}

	var bits bytes.Buffer
	for i, v := range blogs {
		if i-1 > amt {
			break
		}
		components.PostMini(v).Render(context.Background(), &bits)
	}

	return templ.Raw(bits.String()), nil
}
