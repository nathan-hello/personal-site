package main

import (
	_ "embed"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"slices"

	"github.com/nathan-hello/personal-site/db"
	"github.com/nathan-hello/personal-site/render"
	"github.com/nathan-hello/personal-site/router"
	"github.com/nathan-hello/personal-site/utils"
)

const OUTPUT_DIR = "./dist"

const INPUT_BLOG = "./public/content/blog"
const INPUT_PAGES = "./pages"
const INPUT_PUBLIC = "./public"

const DATABASE_URI = "./data.db"

func main() {
	_, err := db.InitDb(DATABASE_URI)
	if err != nil {
		log.Fatal(err)
	}

	err = render.PagesHtml(INPUT_PAGES, OUTPUT_DIR)
	if err != nil {
		log.Fatal(err)
	}
	err = render.Public(INPUT_PUBLIC, OUTPUT_DIR)
	if err != nil {
		log.Fatal(err)
	}

	blogs, err := render.Blogs(INPUT_BLOG, OUTPUT_DIR, true)
	if err != nil {
		log.Fatal(err)
	}

	err = render.Rss(blogs, OUTPUT_DIR)
	if err != nil {
		log.Fatal(err)
	}

	// Currently no static templs, but we could!
	err = render.PagesTempl(OUTPUT_DIR, []render.TemplStaticPages{})
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	for _, v := range router.ApiRoutes {
		mux.Handle(v.Route, v.Middlewares.ThenFunc(v.Hfunc))
	}

	// If dev server, serve /dist as a FileServer
	// If prod, 404 on things that don't match a known route
	// Nginx is responsible for handling static routes without .html
	// E.g. /tv instead of /tv.html
	if slices.Contains(os.Args, "--dev") {
		mux.Handle("/", http.FileServer(http.Dir(OUTPUT_DIR)))
	} else {
		// TODO(nate): what?
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/" {
				http.Redirect(w, r, utils.StatusCodes[404], http.StatusMovedPermanently)
				return
			}
			http.ServeFile(w, r, filepath.Join(OUTPUT_DIR, "index.html"))
		})
	}

	log.Println("Starting webserver on :3000")
	log.Fatal(http.ListenAndServe(":3000", mux))
}
