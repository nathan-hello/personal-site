package main

import (
	"log"
	"os"
	"slices"

	"github.com/nathan-hello/personal-site/db"
	"github.com/nathan-hello/personal-site/render"
	"github.com/nathan-hello/personal-site/router"
)

const INPUT_BLOG = "./public/content/blog"
const INPUT_PAGES = "./pages"
const INPUT_PUBLIC = "./public"

const OUTPUT_PUBLIC = "/var/www/reluekiss.com/public"
const OUTPUT_PRIVATE = "/var/www/reluekiss.com/private"
const OUTPUT_STATIC_FILES = "/var/www/reluekiss.com/public"
const FILE_DATABASE = "file:/var/www/reluekiss.com/private/data.db"
const FILE_CERT = "/var/www/reluekiss.com/private/reluekiss.cert"
const FILE_KEY = "/var/www/reluekiss.com/private/reluekiss.key"

func main() {
	initFiles()
	generate()

	if slices.Contains(os.Args, "--build-only") {
		return
	}

	err := router.SiteRouter(FILE_CERT, FILE_KEY, OUTPUT_PUBLIC)
	if err != nil {
		log.Fatal(err)
	}
}

func initFiles() {
	err := os.RemoveAll(OUTPUT_PUBLIC)
	if err != nil {
		log.Fatal(err)
	}

	err = os.MkdirAll(OUTPUT_PUBLIC, 0544)
	if err != nil {
		log.Fatal(err)
	}

	err = os.MkdirAll(OUTPUT_PRIVATE, 0600)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.InitDb(FILE_DATABASE)
	if err != nil {
		log.Fatal(err)
	}

}

func generate() {
	err := render.Public(INPUT_PUBLIC, OUTPUT_PUBLIC)
	if err != nil {
		log.Fatal(err)
	}

	err = render.PagesHtml(INPUT_PAGES, OUTPUT_PUBLIC)
	if err != nil {
		log.Fatal(err)
	}

	_, err = render.Blogs(INPUT_BLOG, OUTPUT_PUBLIC, true)
	if err != nil {
		log.Fatal(err)
	}

	// Currently no static templs, but we could!
	err = render.PagesTempl(OUTPUT_PUBLIC, []render.TemplStaticPages{})
	if err != nil {
		log.Fatal(err)
	}

}
