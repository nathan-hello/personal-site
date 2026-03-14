package routes

import (
	"context"
	"net/http"

	"github.com/a-h/templ"
	"github.com/nathan-hello/personal-site/src/components"
	"github.com/nathan-hello/personal-site/src/components/layouts"
	"github.com/nathan-hello/personal-site/src/utils"
)

func TvRoute(w http.ResponseWriter, r *http.Request) {
	isStreaming := false;


	child := templ.WithChildren(context.Background(), components.TvPage(true))
	layouts.BaseLayout(components.Header(utils.AsciiTv), components.Meta("NatTV", "Live streaming", "https://reluekiss.com/tv.html", nil )).Render(child, w)
}
