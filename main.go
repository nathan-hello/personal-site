package main

import (
	"log"
	"net/http"

	"github.com/nathan-hello/personal-site/db"
	"github.com/nathan-hello/personal-site/pages"
	"github.com/nathan-hello/personal-site/router"
)

func main() {
	err := router.RenderStaticFiles()
	if err != nil {
		log.Fatal(err)
	}
	err = router.RenderStaticHtml()
	if err != nil {
		log.Fatal(err)
	}
        
        templs := []router.TemplStaticPages{
                {Templ: pages.Index([]db.Blog{}), Route: "/index.html"},
        }

	err = router.RenderStaticTempls(templs)
	if err != nil {
		log.Fatal(err)
	}

	fs := http.FileServer(http.Dir("./dist"))
	http.Handle("/", fs)

	log.Print("Listening on :3000...")
	err = http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatal(err)
	}

}
