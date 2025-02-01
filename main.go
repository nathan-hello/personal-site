package main

import (
	"log"
	"os"
	"slices"

	"github.com/nathan-hello/personal-site/db"
	"github.com/nathan-hello/personal-site/render"
	"github.com/nathan-hello/personal-site/router"
	"github.com/nathan-hello/personal-site/utils"
)

func main() {
        _,err := db.InitDb()
	if err != nil {
		log.Fatal(err)
	}

	err = render.Public()
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

	// Currently no static templs, but we could!
	err = render.PagesTempl([]render.TemplStaticPages{})
	if err != nil {
		log.Fatal(err)
	}

	if slices.Contains(os.Args, "--serve") {
		err = router.SiteRouter()
		if err != nil {
			log.Fatal(err)
		}
	}
}
