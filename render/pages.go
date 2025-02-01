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

func PagesHtml() error {

	err := filepath.Walk(utils.DIR_PAGES, func(path string, info fs.FileInfo, err error) error {
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
                if info.Name() != "index.html" {
                route = strings.TrimSuffix(route,filepath.Ext(info.Name()))
                }
		dist := "dist" + route

		err = writeHtmlFile(f, dist)
		if err != nil {
			return nil
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
	meta.dist = dist

	comp := choosePageLayout(meta)

        f.Seek(0,0)

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
	os.MkdirAll(folder, 0777)
	os.WriteFile(dist, bits.Bytes(), 0777)
	return nil
}

type metadata struct {
	ascii          string
	title          string
	description    string
	overrideLayout string
	dist           string
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
	if strings.Contains(meta.dist, "natalie") {
		layout = "natalie"
	}

	url := strings.TrimPrefix(meta.dist, "dist")
	return registeredPageLayouts[layout](components.Header(meta.ascii), components.Meta(meta.title, meta.description, url, nil))
}
