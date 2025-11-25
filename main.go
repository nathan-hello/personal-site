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

const OUTPUT_PUBLIC = "./dist/public"

const INPUT_BLOG = "./public/content/blog"
const INPUT_PAGES = "./pages"
const INPUT_PUBLIC = "./public"

//go:embed .env
var dotenv string

func main() {
	build := slices.Contains(os.Args, "--build")
	serve := slices.Contains(os.Args, "--serve")
	isDev := slices.Contains(os.Args, "--dev")

	if isDev {
		build = true
		serve = true
	}

	err := utils.ParseDotenv(dotenv)
	if err != nil {
		log.Fatal(err)
	}

	log.Print("before db init")
	_, err = db.InitDb(utils.Env().DATABASE_URI)
	if err != nil {
		log.Fatal(err)
	}
	log.Print("after db init")

	if !build && !serve {
		log.Fatal("neither --build or --serve was given: choose one!")
	}

	if build {
		generate()
	}

	if serve {
		serveHttp()
	}
}

func generate() {
	err := render.PagesHtml(INPUT_PAGES, OUTPUT_PUBLIC)
	if err != nil {
		log.Fatal(err)
	}
	err = render.Public(INPUT_PUBLIC, OUTPUT_PUBLIC)
	if err != nil {
		log.Fatal(err)
	}

	blogs, err := render.Blogs(INPUT_BLOG, OUTPUT_PUBLIC, true)
	if err != nil {
		log.Fatal(err)
	}

	err = render.Rss(blogs, OUTPUT_PUBLIC)
	if err != nil {
		log.Fatal(err)
	}

	// Currently no static templs, but we could!
	err = render.PagesTempl(OUTPUT_PUBLIC, []render.TemplStaticPages{})
	if err != nil {
		log.Fatal(err)
	}

}

func serveHttp() {
	mux := http.NewServeMux()
	for _, v := range router.ApiRoutes {
		mux.Handle(v.Route, v.Middlewares.ThenFunc(v.Hfunc))
	}

	// If dev server, serve /dist as a FileServer
	// If prod, 404 on things that don't match a known route
	// Nginx is responsible for handling static routes without .html
	// E.g. /tv instead of /tv.html
	if slices.Contains(os.Args, "--dev") {
		mux.Handle("/", http.FileServer(http.Dir(OUTPUT_PUBLIC)))
	} else {
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/" {
				http.Redirect(w, r, utils.StatusCodes[404], http.StatusMovedPermanently)
				return
			}
			http.ServeFile(w, r, filepath.Join(OUTPUT_PUBLIC, "index.html"))
		})
	}

	mux.HandleFunc("/git-hook", utils.HookHandler)

	log.Println("Starting webserver on :3000")
	log.Fatal(http.ListenAndServe(":3000", mux))
}
