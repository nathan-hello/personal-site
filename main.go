package main

import (
	"log"
	"net/http"
	"os"
	"slices"
	"strings"

	"github.com/nathan-hello/personal-site/render"
	"github.com/nathan-hello/personal-site/utils"
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

	blogs, err := render.Blogs(true)
	if err != nil {
		log.Fatal(err)
	}

	slices.SortFunc(blogs, func(a, b utils.Blog) int {
		return b.Frnt.Date.Compare(a.Frnt.Date)
	})
	err = render.PagesTempl([]render.TemplStaticPages{
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
