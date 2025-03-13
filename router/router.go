package router

import (
	"fmt"
	"net/http"

	"github.com/justinas/alice"
	"github.com/nathan-hello/personal-site/router/routes"
)

type Site struct {
	route       string
	hfunc       http.HandlerFunc
	middlewares alice.Chain
}

func SiteRouter(filesDir string) error {
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

	fmt.Printf("Listening on :3001...")
	err := http.ListenAndServe(":3001", nil)
	if err != nil {
		return err
	}

	return nil
}
