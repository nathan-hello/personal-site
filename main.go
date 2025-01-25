package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/nathan-hello/personal-site/pages"
	"github.com/nathan-hello/personal-site/render"
)

func main() {
	err := render.Public()
	if err != nil {
		log.Fatal(err)
	}
	err = render.PagesHtml()
	if err != nil {
		log.Fatal(err)
	}

        blogs,err := render.Blogs()
	if err != nil {
		log.Fatal(err)
	}

	err = render.PagesTempl([]render.TemplStaticPages{
		{Templ: pages.Index(blogs), Route: "/index.html"},
	})
	if err != nil {
		log.Fatal(err)
	}

	serve()
}

func serve() {
	serve := false
	for _, v := range os.Args {
		if strings.Contains(v, "--serve") {
			serve = true
		}
	}

	if serve {
		fs := http.FileServer(http.Dir("./dist"))
		http.Handle("/", fs)

		log.Print("Listening on :3000...")
		err := http.ListenAndServe(":3000", nil)
		if err != nil {
			log.Fatal(err)
		}
	}

}
