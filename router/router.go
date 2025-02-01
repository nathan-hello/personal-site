package router

import (
	"net/http"

	"github.com/justinas/alice"
	"github.com/nathan-hello/personal-site/router/routes"
)

func SiteRouter() error {

	type Site struct {
		route       string
		hfunc       http.HandlerFunc
		middlewares alice.Chain
	}

	sites := []Site{
		{route: "/api/comments/{id}",
			hfunc: routes.ApiComments,
			middlewares: alice.New(
				RejectSubroute("/"),
				Logging,
				AllowMethods("GET"),
			)},
	}

	for _, v := range sites {
		http.Handle(v.route, v.middlewares.ThenFunc(v.hfunc))
	}

	fs := http.FileServer(http.Dir("./dist"))
	http.Handle("/", fs)

        err := http.ListenAndServe(":3000", nil)
		if err != nil {
                return err
		}
        return nil
}

