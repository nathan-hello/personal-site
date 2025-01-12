package router

import (
	"net/http"

	"github.com/justinas/alice"
	"github.com/nathan-hello/personal-site/router/routes"
)

func SiteRouter() {

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

	http.ListenAndServe(":3000", nil)
}

