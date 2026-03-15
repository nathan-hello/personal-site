package router

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/justinas/alice"
	"github.com/nathan-hello/personal-site/src/router/routes"
	"github.com/nathan-hello/personal-site/src/router/routes/chat"
)

type Site struct {
	Route       string
	Hfunc       http.HandlerFunc
	Middlewares alice.Chain
}

var ApiRoutes = []Site{
	{Route: "/api/comments/{id}",
		Hfunc: routes.ApiComments,
		Middlewares: alice.New(
			Logging,
			AllowMethods("GET", "POST"),
		)},
	{Route: "/tv",
		Hfunc: routes.TvRoute,
		Middlewares: alice.New(
			Logging,
			AllowMethods("GET", "POST"),
		)},
	{Route: "/tv.html",
		Hfunc: routes.TvRoute,
		Middlewares: alice.New(
			Logging,
			AllowMethods("GET", "POST"),
		)},
	// {Route: "/api/comment-delete",
	// 	Hfunc: routes.ApiCommentsDelete,
	// 	Middlewares: alice.New(
	// 		Logging,
	// 		AllowMethods("POST"),
	// 	),
	// },
	{Route: "/api/captcha",
		Hfunc: routes.ApiCaptcha,
		Middlewares: alice.New(
			Logging,
			AllowMethods("GET", "POST"),
		),
	},
	{Route: "/i/{id}",
		Hfunc: routes.ApiCommentImage,
		Middlewares: alice.New(
			Logging,
			AllowMethods("GET"),
		),
	},
	{Route: "/api/chat",
		Hfunc: chat.ApiChat,
		Middlewares: alice.New(
			Logging,
			AllowMethods("GET", "POST"),
		),
	},
	{Route: "/chat",
		Hfunc: chat.Chat,
		Middlewares: alice.New(
			Logging,
			AllowMethods("GET"),
		),
	},
	{Route: "/weather",
		Hfunc: routes.Weather,
		Middlewares: alice.New(
			Logging,
			AllowMethods("GET"),
		),
	},
	{Route: "/weather/{location}",
		Hfunc: routes.Weather,
		Middlewares: alice.New(
			Logging,
			AllowMethods("GET"),
		),
	},
	{Route: "/ws/v1/chat/html",
		Hfunc: chat.ChatSocket,
		Middlewares: alice.New(
			Logging,
		)},
	{Route: "/api/v1/chat/message",
		Hfunc: chat.ApiChat,
		Middlewares: alice.New(
			AllowMethods("POST"),
		)},
}

func RegisterAuth(mux *http.ServeMux) {
	// Redirect handlers for exact paths without trailing slash
	mux.Handle("/auth", http.RedirectHandler("/auth/", http.StatusMovedPermanently))
	mux.Handle("/api/auth", http.RedirectHandler("/api/auth/", http.StatusMovedPermanently))

	// Proxy handler for paths with trailing slash
	proxy := httputil.NewSingleHostReverseProxy(&url.URL{
		Scheme: "http",
		Host:   "localhost:3005",
	})
	mux.Handle("/auth/", proxy)
	mux.Handle("/api/auth/", proxy)
}
