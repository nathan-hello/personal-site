package router

import (
	"net/http"

	"github.com/justinas/alice"
	"github.com/nathan-hello/personal-site/router/routes"
)

type Site struct {
	route       string
	hfunc       http.HandlerFunc
	middlewares alice.Chain
}

var ApiRoutes = []Site{
	{route: "/api/comments/{id}",
		hfunc: routes.ApiComments,
		middlewares: alice.New(
			Logging,
			AllowMethods("GET", "POST"),
		)},
    {route: "/api/captcha",
        hfunc: routes.ApiCaptcha,
        middlewares: alice.New(
            Logging,
            AllowMethods("GET", "POST"),
        ),
    },
}

func RegisterApiHttpHandler() {
	for _, v := range ApiRoutes {
		http.Handle(v.route, v.middlewares.ThenFunc(v.hfunc))
	}
}
