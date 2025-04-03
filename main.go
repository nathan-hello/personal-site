package main

import (
	_ "embed"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"

	"github.com/fsnotify/fsnotify"

	"github.com/nathan-hello/personal-site/db"
	"github.com/nathan-hello/personal-site/render"
	"github.com/nathan-hello/personal-site/router"
	"github.com/nathan-hello/personal-site/utils"
)

const OUTPUT_PUBLIC = "./dist/public"
const OUTPUT_PRIVATE = "./dist/private"

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

	_, err = db.InitDb(utils.Env().DATABASE_URI)
	if err != nil {
		log.Fatal(err)
	}

	if !build && !serve {
		log.Fatal("neither --build or --serve was given: choose one!")
	}

	if build {
		generate()
	}

	if serve {
		go watchFiles()
		go serveHttp()
		select {}
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
	router.RegisterApiHttpHandler()

	if slices.Contains(os.Args, "--dev") {
        router.Index(true, OUTPUT_PUBLIC)
	}

	fmt.Printf("Listening on port :3000 for routes: %v\n", router.ApiRoutes)

	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatal(err)
	}

}

func watchFiles() {
    watcher, err := fsnotify.NewWatcher()
    if err != nil {
        log.Fatal(err)
    }
    defer watcher.Close()

	err = filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        if info.IsDir() {
            if info.Name() == "dist" {
                return filepath.SkipDir
            }
            return watcher.Add(path)
        }
        if strings.Contains(info.Name(), "_templ") {
            return nil
        }
        return nil
    })
    if err != nil {
        log.Fatal(err)
    }

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}

			if event.Op&(fsnotify.Write|fsnotify.Create) != 0 {
				ext := filepath.Ext(event.Name)
				if ext == ".mdx" {
					make("build/css")
					generate()
				}
				if ext == ".go" || ext == ".sql" || ext == ".templ"  {
					// TODO: why does this loop?
					// make("build/css")
					// make("build/templ")
					// make("build/sqlc")
					// TODO: make project rebuild and run itself
				}
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Println("error:", err)
		}
	}
}

func make(arg string) {
	cmd := exec.Command("make", arg)
	cmd.Stdout = io.Discard
	cmd.Stderr = io.Discard
	fmt.Printf("INFO: rebuild called make %s\n", arg)
	if err := cmd.Start(); err != nil {
        log.Printf("Rebuild failed: %v", err)
        return
    }
}
