package routes

import (
	"net/http"

	"github.com/nathan-hello/personal-site/auth"
	"github.com/nathan-hello/personal-site/components"
	"github.com/nathan-hello/personal-site/router/routes/chat"
)

func SignIn(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		components.SignIn(auth.AuthResult{}, auth.SignInData{}).Render(r.Context(), w)
		return
	}

	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	res := auth.SignIn(auth.SignInData{
		UserOrEmail: r.FormValue("user"),
		Password:    r.FormValue("password"),
	})

	if len(res.Errors) > 0 {
		components.SignIn(res, auth.SignInData{
			UserOrEmail: r.FormValue("user"),
			Password:    r.FormValue("password"),
		}).Render(r.Context(), w)
		return
	}

	access, refresh, err := auth.NewTokenPair(&auth.JwtParams{
		UserId:   res.User.ID,
		Username: res.User.Username,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	auth.SetTokenCookies(w, access, refresh)
	chat.Chat(w, r)
}

func SignOut(w http.ResponseWriter, r *http.Request) {
	auth.DeleteJwtCookies(w)
	components.SignIn(auth.AuthResult{}, auth.SignInData{}).Render(r.Context(), w)
}
