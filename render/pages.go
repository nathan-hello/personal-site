package render

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/a-h/templ"
	"github.com/nathan-hello/personal-site/components"
	"github.com/nathan-hello/personal-site/layouts"
	"github.com/nathan-hello/personal-site/utils"
)

// This renders .html files in the pages/ dir.

func PagesHtml(input, output string) error {

	err := filepath.Walk(input, func(path string, info fs.FileInfo, err error) error {
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

		route := strings.TrimPrefix(path, "pages")

                // Generated files keep their extension.
                // Tell nginx to try without .html and look for the same path but with .html
		// if info.Name() != "index.html" {
		// 	route = strings.TrimSuffix(route, filepath.Ext(info.Name()))
		// }

		dist := output + route

		err = writeHtmlFile(f, dist)
		if err != nil {
			return err
		}

		return nil

	})
	if err != nil {
		return err
	}

	return nil
}

func writeHtmlFile(f *os.File, dist string) error {
	meta := parsePagesFrontmatter(f)
	meta.path = dist
        meta.url = strings.TrimSuffix(meta.path, ".html")
	meta.url = strings.TrimPrefix(meta.path, "dist")

	comp := choosePageLayout(meta)

	f.Seek(0, 0)

	content, err := io.ReadAll(f)
	if err != nil {
		return err
	}

	renderedFile, err := RenderCustomComponents(string(content))
	if err != nil {
		return err
	}

	var bits bytes.Buffer
	childrenCtx := templ.WithChildren(context.Background(), templ.Raw(renderedFile))
	err = comp.Render(childrenCtx, &bits)
	if err != nil {
		return err
	}

	parts := strings.Split(dist, "/")
	folder := strings.Join(parts[:len(parts)-1], "/")
	fmt.Printf("INFO: writing file %s in folder %s\n", dist, folder)
	err = os.MkdirAll(folder, 0777)
	if err != nil {
		return err
	}
	err = os.WriteFile(dist, bits.Bytes(), 0777)
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
        url string
}

func parsePagesFrontmatter(f *os.File) metadata {

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
		meta.overrideLayout = m[1]
	}

	if scanner.Err() != nil {
		fmt.Printf("WARN: metadata could not be rendered for file at %s. err: %s\n", f.Name(), scanner.Err().Error())
		meta = metadata{ascii: utils.AsciiNat_e, title: "reluekiss.com", description: "Nat/e. We are Boingus."}
	}

	if meta.ascii == "" {
		meta.ascii = utils.AsciiNat_e
	}
	if meta.title == "" {
		meta.title = "Nat/e - reluekiss.com"
	}
	if meta.description == "" {
		meta.description = "We are boingus."
	}

	return meta
}

var layoutMap map[string]layouts.LayoutComponent
var registeredPageLayouts = map[string]layouts.LayoutComponent{
	"natalie": layouts.NatalieFullPage,
	"default": layouts.BaseLayout,
}

func choosePageLayout(meta metadata) templ.Component {
	layout := "default"
	if strings.Contains(meta.url, "natalie") {
		layout = "natalie"
	}

	return registeredPageLayouts[layout](components.Header(meta.ascii), components.Meta(meta.title, meta.description, meta.url, nil))
}
