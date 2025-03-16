package render

import (
	"bytes"
	"context"
	"fmt"
	"os"
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
		components.PostMini(v).Render(context.Background(), &bits)
		if i > amt {
			break
		}
	}

	return templ.Raw(bits.String()), nil
}

var characterSheetLabels = map[string]string{
        "txt": "Character Sheet for",
        "diff": "Updated Character Sheet for",
}

// <CharacterSheet story="cyberpunk" character="natasha" version="1.1" type="diff | txt" label?="string with periods instead of spaces, with CHARACTER always following immediately after"/>
func characterSheet(c component) (templ.Component, error) {

	story, ok := c.Attributes["story"]
	if !ok {
		return nil, fmt.Errorf("component did not give required attribute \"story\": %#v", c)
	}
	character, ok := c.Attributes["character"]
	if !ok {
		return nil, fmt.Errorf("component did not give required attribute \"character\": %#v", c)
	}
	version, ok := c.Attributes["version"]
	if !ok {
		return nil, fmt.Errorf("component did not give required attribute \"version\": %#v", c)
	}

	ext, ok := c.Attributes["type"]
	if !ok {
		return nil, fmt.Errorf("component did not give required attribute \"type\": %#v", c)
	}

	label, ok := c.Attributes["label"]
	if !ok {
                defaultLabel, ok := characterSheetLabels[ext]
                if ok {
                        label = defaultLabel
                } else {
                        label = strings.ReplaceAll(label, ".", " ")
                }
	}

	path := fmt.Sprintf("./public/content/character-sheets/%s/%s-%s.%s", story, character, version, ext)

	fmt.Printf("component: %#v\n, path: %s\n", c, path)

	f, err := os.ReadFile(path)

	if err != nil {
		return nil, fmt.Errorf("err reading file for component %#v, os.ReadFile err: %s", c, err)
	}

	return components.CharacterSheet(label, character, version, string(f)), nil

}
