package router

import (
	"net/http"

	"github.com/justinas/alice"
	"github.com/nathan-hello/personal-site/router/routes"
	"github.com/nathan-hello/personal-site/router/routes/chat"
)

type Site struct {
	route       string
	hfunc       http.HandlerFunc
	middlewares alice.Chain
}

var ApiRoutes = []Site{
	{route: "/",
		hfunc: Index(false, ""),
		middlewares: alice.New(
			Logging,
		),
	},
	{route: "/api/comments/{id}",
		hfunc: routes.ApiComments,
		middlewares: alice.New(
			Logging,
			AllowMethods("GET", "POST"),
		)},
	{route: "/api/comment-delete",
		hfunc: routes.ApiCommentsDelete,
		middlewares: alice.New(
			Logging,
			AllowMethods("POST"),
		),
	},
	{route: "/api/captcha",
		hfunc: routes.ApiCaptcha,
		middlewares: alice.New(
			Logging,
			AllowMethods("GET", "POST"),
		),
	},
	{route: "/i/{id}",
		hfunc: routes.ApiCommentImage,
		middlewares: alice.New(
			Logging,
			AllowMethods("GET"),
		),
	},
	{route: "/api/chat",
		hfunc: chat.ApiChat,
		middlewares: alice.New(
			Logging,
			AllowMethods("GET", "POST"),
			InjectClaimsOnValidToken,
		),
	},
	{route: "/bear/chat",
		hfunc: chat.BearChat,
		middlewares: alice.New(
			Logging,
			AllowMethods("GET", "POST"),
			InjectClaimsOnValidToken,
		),
	},
	{route: "/bear/login",
		hfunc: routes.BearLogin,
		middlewares: alice.New(
			Logging,
			AllowMethods("GET", "POST"),
			InjectClaimsOnValidToken,
		),
	},
	{route: "/bear/signout",
		hfunc: routes.BearSignOut,
		middlewares: alice.New(
			Logging,
			AllowMethods("GET", "POST"),
			InjectClaimsOnValidToken,
		),
	},
	{route: "/login",
		hfunc: routes.Login,
		middlewares: alice.New(
			Logging,
			AllowMethods("GET", "POST"),
			InjectClaimsOnValidToken,
		),
	},
	{route: "/signout",
		hfunc: routes.SignOut,
		middlewares: alice.New(
			Logging,
			AllowMethods("GET", "POST"),
			InjectClaimsOnValidToken,
		),
	},
}

func RegisterApiHttpHandler() {
	for _, v := range ApiRoutes {
		http.Handle(v.route, v.middlewares.ThenFunc(v.hfunc))
	}
}
