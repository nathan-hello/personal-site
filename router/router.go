package router

import (
	"net/http"
	"os"
	"slices"

	"github.com/justinas/alice"
	"github.com/nathan-hello/personal-site/router/routes"
)

type Site struct {
	route       string
	hfunc       http.HandlerFunc
	middlewares alice.Chain
}

func SiteRouter(cert, key, filesDir string) error {
	sites := []Site{
		{route: "/api/comments/{id}",
			hfunc: routes.ApiComments,
			middlewares: alice.New(
				Logging,
				AllowMethods("GET", "POST"),
			)},
	}

	for _, v := range sites {
		http.Handle(v.route, v.middlewares.ThenFunc(v.hfunc))
	}

	fs := http.FileServer(http.Dir(filesDir))
	http.Handle("/", fs)

	if slices.Contains(os.Args, "--prod-server") {
		go func() {
			http.ListenAndServe(":80", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				http.Redirect(w, r, "https://"+r.Host+r.RequestURI, http.StatusMovedPermanently)
			}))
		}()

		err := http.ListenAndServeTLS(":443", cert, key, nil)
		if err != nil {
			return err
		}
	} else {
		err := http.ListenAndServe(":3000", nil)
		if err != nil {
			return err
		}
	}

	return nil
}
