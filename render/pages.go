package render

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/a-h/templ"
	"github.com/nathan-hello/personal-site/components"
	"github.com/nathan-hello/personal-site/layouts"
	"github.com/nathan-hello/personal-site/render/customs"
	"github.com/nathan-hello/personal-site/utils"
)

// This renders .html files in the pages/ dir.

func PagesHtml() error {

	err := filepath.Walk("pages", func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		if filepath.Ext(info.Name()) != ".html" {
			return nil
		}

		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		meta := parseHtmlHeader(f)
		comp := chooseLayout(meta)

		renderedFile, err := customs.RenderCustomComponents(f)
		if err != nil {
			return err
		}
		fmt.Println(renderedFile)

		route := strings.TrimPrefix(path, "pages")
		dist := "dist" + route

		var bits bytes.Buffer
		childrenCtx := templ.WithChildren(context.Background(), templ.Raw(renderedFile))
		err = comp.Render(childrenCtx, &bits)
		if err != nil {
			return err
		}

		folder := strings.TrimSuffix(dist, info.Name())
		os.MkdirAll(folder, 0777)
		os.WriteFile(dist, bits.Bytes(), 0777)

		return nil

	})
	if err != nil {
		return err
	}

	return nil
}

type metadata struct {
	ascii          string
	title          string
	description    string
	overrideLayout string
	path           string
}

func parseHtmlHeader(f *os.File) metadata {

	var sb strings.Builder
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		sb.WriteString(line + "\n")
		if strings.Contains(line, "-->") {
			break
		}
	}

	asciiRe := regexp.MustCompile(`(?s)<ascii>(.*?)</ascii>`)
	titleRe := regexp.MustCompile(`<title>(.*?)</title>`)
	descRe := regexp.MustCompile(`<description>(.*?)</description>`)
	layoRe := regexp.MustCompile(`<layout>(.*?)</layout>`)

	content := sb.String()

	meta := metadata{}
	if m := asciiRe.FindStringSubmatch(content); len(m) > 1 {
		meta.ascii = m[1]
	}
	if m := titleRe.FindStringSubmatch(content); len(m) > 1 {
		meta.title = m[1]
	}
	if m := descRe.FindStringSubmatch(content); len(m) > 1 {
		meta.description = m[1]
	}
	if m := layoRe.FindStringSubmatch(content); len(m) > 1 {
		meta.description = m[1]
	}

	if scanner.Err() != nil {
		fmt.Printf("WARN: metadata could not be rendered for file at %s. err: %s\n", f.Name(), scanner.Err().Error())
		meta = metadata{ascii: utils.AsciiNat_e, title: "reluekiss.com", description: "Nat/e. We are Boingus."}
	}

	return meta
}

var layoutMap map[string]layouts.LayoutComponent
var registeredLayouts = map[string]layouts.LayoutComponent{
	"natalie": layouts.NatalieFullPage,
	"default": layouts.IndexLayout,
}

func chooseLayout(meta metadata) templ.Component {
	layout := "default"
	if strings.Contains(meta.path, "natalie") {
		layout = "natalie"
	}
	url := strings.Split(meta.path, "pages")[1]

	return registeredLayouts[layout](components.Header(meta.title), components.Meta(meta.title, meta.description, url, components.DefaultImage))
}
