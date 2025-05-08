package router

import (
	"net/http"

	"github.com/justinas/alice"
	"github.com/nathan-hello/personal-site/router/routes"
	"github.com/nathan-hello/personal-site/router/routes/chat"
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
	{Route: "/api/comment-delete",
		Hfunc: routes.ApiCommentsDelete,
		Middlewares: alice.New(
			Logging,
			AllowMethods("POST"),
		),
	},
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
			InjectClaimsOnValidToken,
		),
	},
	{Route: "/signin",
		Hfunc: routes.SignIn,
		Middlewares: alice.New(
			Logging,
			AllowMethods("GET", "POST"),
			InjectClaimsOnValidToken,
		),
	},
	{Route: "/auth/signout",
		Hfunc: routes.SignOut,
		Middlewares: alice.New(
			Logging,
			AllowMethods("GET", "POST"),
			InjectClaimsOnValidToken,
		),
	},
	{Route: "/chat",
		Hfunc: chat.Chat,
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
}

func RegisterApiHttpHandler() {
	for _, v := range ApiRoutes {
		http.Handle(v.Route, v.Middlewares.ThenFunc(v.Hfunc))
	}
}
