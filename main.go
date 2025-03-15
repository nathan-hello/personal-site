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
	"public":  "./dist/public",
	"private": "./dist/private",
	"db":      "/var/www/reluekiss.com/private/data.db",
}

var dev map[string]string = map[string]string{
	"public":  "./dist/public",
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

	build := slices.Contains(os.Args, "--build")
	serve := slices.Contains(os.Args, "--serve")
	if build && serve {
		log.Fatal("both --build and --serve was given: choose one!")
	}
	if !build && !serve {
		log.Fatal("neither --build or --serve was given: choose one!")
	}

	if build {
		generate(m)
	}

	if serve {
		startHttp(m)
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

func startHttp(m map[string]string) {
	_, err := db.InitDb(m["db"])
	if err != nil {
		log.Fatal(err)
	}
	router.RegisterApiHttpHandler()

	if slices.Contains(os.Args, "--dev") {
		http.Handle("/", http.FileServer(http.Dir(m["public"])))
	}

	http.ListenAndServe(":3000", nil)

}
