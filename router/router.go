package router

import (
	"io"
	"log"
	"net/http"
	"os"

	"github.com/justinas/alice"
	natauth "github.com/nathan-hello/nat-auth"
	"github.com/nathan-hello/nat-auth/storage"
	"github.com/nathan-hello/nat-auth/ui"
	"github.com/nathan-hello/nat-auth/web"
	"github.com/nathan-hello/personal-site/router/routes"
	"github.com/nathan-hello/personal-site/router/routes/chat"
)

type Site struct {
	Route       string
	Hfunc       http.HandlerFunc
	Middlewares alice.Chain
}
func auth() natauth.Handlers {
	store, err := storage.NewValkey("127.0.0.1:6379", "app")
	if err != nil {
		log.Printf("error in NewValkey: %s\n", err.Error())
		log.Fatal(err)
	}

	handlers, err := natauth.New(natauth.Params{
		JwtConfig:  web.PasswordJwtParams{Secret: "secret"},
		Storage:    store,
		LogWriters: []io.Writer{os.Stdout},
		Theme: ui.Theme{
			Primary: ui.ColorScheme{
				Light: "#262626",
				Dark:  "#262626",
			},
			Background: ui.ColorScheme{
				Light: "#171717",
				Dark:  "#171717",
			},
			Logo: ui.ColorScheme{
				Light: "https://reluekiss.com/favicon.svg",
				Dark:  "https://reluekiss.com/favicon.svg",
			},
			Title:   "Nat/e",
			Favicon: "https://reluekiss.com/favicon.svg",
			Radius:  "none",
			Font: ui.Font{
				Family: "Varela Round, sans-serif",
				Scale:  "1",
			},
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	return handlers
}

var authHandler = auth()


var ApiRoutes = []Site{
	{Route: "/api/comments/{id}",
		Hfunc: routes.ApiComments,
		Middlewares: alice.New(
			Logging,
			AllowMethods("GET", "POST"),
			authHandler.Middleware,
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
}

func RegisterApiHttpHandler() {
	for _, v := range ApiRoutes {
		http.Handle(v.Route, v.Middlewares.ThenFunc(v.Hfunc))
	}
}
