package main

import (
	"log"
	"net/http"
	"os"
	"slices"

	"github.com/nathan-hello/personal-site/db"
	"github.com/nathan-hello/personal-site/render"
	"github.com/nathan-hello/personal-site/router"
)

var prod map[string]string = map[string]string{
	"public":  "/var/www/reluekiss.com/public",
	"private": "/var/www/reluekiss.com/private",
	"db":      "/var/www/reluekiss.com/private/data.db",
}

var dev map[string]string = map[string]string{
	"public":  "./dist",
	"private": "./dist/private",
	"db":      ":memory:",
}

const INPUT_BLOG = "./public/content/blog"
const INPUT_PAGES = "./pages"
const INPUT_PUBLIC = "./public"

func main() {
	m := prod
	if slices.Contains(os.Args, "--dev") {
		m = dev
	}

	initFiles(m)
	generate(m)

	if slices.Contains(os.Args, "--build-only") {
		return
	}

	err := router.SiteRouter(m["public"])
	if err != nil {
		log.Fatal(err)
	}

	if slices.Contains(os.Args, "--dev") {
		fs := http.FileServer(http.Dir("./dir"))
		http.Handle("/", fs)
		http.ListenAndServe(":8000", nil)
	}
}

func initFiles(m map[string]string) {
	err := os.RemoveAll(m["public"])
	if err != nil {
		log.Fatal(err)
	}

	err = os.MkdirAll(m["public"], 0744)
	if err != nil {
		log.Fatal(err)
	}

	err = os.MkdirAll(m["private"], 0700)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.InitDb(m["db"])
	if err != nil {
		log.Fatal(err)
	}

}

func generate(m map[string]string) {

	err := render.PagesHtml(INPUT_PAGES, m["public"])
	if err != nil {
		log.Fatal(err)
	}
	err = render.Public(INPUT_PUBLIC, m["public"])
	if err != nil {
		log.Fatal(err)
	}

	_, err = render.Blogs(INPUT_BLOG, m["public"], true)
	if err != nil {
		log.Fatal(err)
	}

	// Currently no static templs, but we could!
	err = render.PagesTempl(m["public"], []render.TemplStaticPages{})
	if err != nil {
		log.Fatal(err)
	}

}
