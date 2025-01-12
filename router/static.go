package router

import (
	"bytes"
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/a-h/templ"
	"github.com/nathan-hello/personal-site/components"
	"github.com/nathan-hello/personal-site/layouts"
	"github.com/nathan-hello/personal-site/utils"
)

type Static struct {
	route       string
	filepath    string
	contentType string
}

func RenderStaticFiles() error {
	files := []Static{}

	err := filepath.Walk("public", func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		ext := filepath.Ext(info.Name())

		route := strings.TrimPrefix(path, "public") // keep "/" in beginning
		var contentType = ""

		if ext == ".js" {
			contentType = "text/javascript"
		}
		if ext == ".css" {
			contentType = "text/css"
		}

		files = append(files, Static{route: route, filepath: path, contentType: contentType})

		dist := "dist" + route

		bits, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		folder := strings.TrimSuffix(dist, info.Name())
		//                fmt.Println(folder)
		os.MkdirAll(folder, 0777)
		os.WriteFile(dist, bits, 0777)

		return err
	})

	if err != nil {
		return err
	}
	if len(files) == 0 {
		return fmt.Errorf("no static files: %#v", files)
	}

	return nil
}

func RenderStaticHtml() error {
	htmls := []Static{}

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

		var comp templ.Component
		var meta = &MetaData{}
		meta, err = ParseHTMLHeader(path)
		if err != nil {
			meta = &MetaData{ASCII: utils.AsciiNat_e, Title: "reluekiss.com", Description: "Nat/e. We are Boingus."}
		}

		if strings.Contains(path, "natalie") {
			comp = layouts.NatalieFullPage(components.Header(meta.ASCII), components.Meta(meta.Title, meta.Description, "", nil))
		} else {
			comp = layouts.IndexLayout(components.Header(meta.ASCII), components.Meta(meta.Title, meta.Description, "", nil))
		}

		route := strings.TrimPrefix(path, "pages")
		dist := "dist" + route

		var bits bytes.Buffer
		err = comp.Render(context.Background(), &bits)
		if err != nil {
			return err
		}

		folder := strings.TrimSuffix(dist, info.Name())
		os.MkdirAll(folder, 0777)
		os.WriteFile(dist, bits.Bytes(), 0777)

		htmls = append(htmls, Static{route: route, filepath: dist, contentType: "text/html"})

		return nil

	})
	if err != nil {
		return err
	}

	return nil
}

type TemplStaticPages struct {
	Templ templ.Component
	Route string
}

func RenderStaticTempls(templs []TemplStaticPages) error {

	var bits bytes.Buffer
	for _, v := range templs {
		v.Templ.Render(context.Background(), &bits)
		parts := strings.Split(v.Route, "/")
		if len(parts) > 1 {
			parts = parts[:len(parts)-1]
		}
		folder := "dist" + strings.Join(parts, "/")
                fmt.Println(folder)
		os.MkdirAll(folder, 0777)
                os.WriteFile("dist" + v.Route, bits.Bytes(), 0777)
	}

	return nil
}

// func StaticRouter(files []Static) error {
// 	for _, v := range files {
// 		// closure shenanigans
// 		file := v.filepath
// 		route := v.route
// 		contentType := v.contentType
// 		middles := alice.New(Logging, AllowMethods("GET"), CreateHeader("Content-Type", v.contentType))
// 		http.Handle(v.route, middles.ThenFunc(func(w http.ResponseWriter, r *http.Request) { http.ServeFile(w, r, file) }))
// 		log.Printf("Creating route: %v, for file: %v, with Content-Type %v\n", route, file, contentType)
// 	}
// 	return nil
//
// }
