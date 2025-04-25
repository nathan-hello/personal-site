package main

import (
	_ "embed"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"syscall"
	"time"

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

	nameCh := make(chan string)
	if serve {
		go watchFiles(nameCh)
		go serveHttp(nameCh)
		nameCh <- "Initial Start"
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

func serveHttp(nameCh chan string) {
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

   fileName := <- nameCh
   fmt.Printf("%v Listening on port :3000 due to change in %v for routes: %v\n", time.Now().Format("15:04:05"), fileName, router.ApiRoutes)
   log.Fatal(http.ListenAndServe(":3000", mux))

}

func watchFiles(nameCh chan string) {
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

    var rebuildTimer *time.Timer
    
    for {
        select {
        case event := <- watcher.Events:
			nameCh <- event.Name
            ext := filepath.Ext(event.Name)
            if ext == ".mdx" || ext == ".html" {
                if rebuildTimer != nil {
                    rebuildTimer.Stop()
                }
                rebuildTimer = time.AfterFunc(200*time.Millisecond, func() {
                    exec.Command("bun", "run", "tailwindcss", "-i", "./public/css/tw-input.css", "-o", "./public/css/tw-output.css").Run()
                    generate()
                })
            }
            if ext == ".go" || ext == ".sql" || ext == ".templ" {
                if rebuildTimer != nil {
                    rebuildTimer.Stop()
                }
                rebuildTimer = time.AfterFunc(200*time.Millisecond, func() {
                    exec.Command("go", "run", "github.com/sqlc-dev/sqlc/cmd/sqlc@v1.27.0", "generate").Run()
                    exec.Command("bun", "run", "tailwindcss", "-i", "./public/css/tw-input.css", "-o", "./public/css/tw-output.css").Run()
                    exec.Command("templ", "generate").Run()
    
                    exe, _ := os.Executable()
                    exec.Command("go", "build", "-o", exe, ".").Run()
    
                    syscall.Exec(exe, append([]string{exe}, os.Args[1:]...), os.Environ())
                })
            }
    
        case err := <-watcher.Errors:
            log.Println("watcher error:", err)
        }
    }
}
