package routes

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/nathan-hello/personal-site/src/auth"
)

func SignUp(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(auth.ClaimsContextKey).(auth.CustomClaims)
	if ok {
		HandleRedirect(w, r, fmt.Sprintf("/profile/%s", claims.Username), nil)
		return
	}

	returnFormWithErrors := func(errs []auth.AuthError) {
		fmt.Printf("%#v\n", errs)
		w.WriteHeader(401)
		// components.SignUpForm(components.RenderAuthError(errs)).Render(r.Context(), w)
		errs = nil
	}

	if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			returnFormWithErrors([]auth.AuthError{
				{Err: auth.ErrParseForm},
			})
		}

		cred := auth.SignUpCredentials{
			Username: r.FormValue("username"),
			Password: r.FormValue("password"),
			PassConf: r.FormValue("password-confirmation"),
			Email:    r.FormValue("email"),
		}

		username, userId, errs := cred.SignUp()

		if errs != nil {
			returnFormWithErrors(errs)
			return
		}

		access, refresh, err := auth.NewTokenPair(
			&auth.JwtParams{
				Username: username,
				UserId:   userId,
				Family:   uuid.New(),
			})

		if err != nil {
			returnFormWithErrors([]auth.AuthError{
				{Err: err},
			})
		}

		auth.SetTokenCookies(w, access, refresh)
		w.Header().Set("HX-Redirect", fmt.Sprintf("/profile/%v", username))
		return

	}

	if r.Method == "GET" {
		w.WriteHeader(200)
		// components.SignUp().Render(r.Context(), w)
		return
	}
}

func SignIn(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(auth.ClaimsContextKey).(auth.CustomClaims)
	if ok {
		HandleRedirect(w, r, fmt.Sprintf("/profile/%s", claims.Username), nil)
		return
	}

	if r.Method == "POST" {
		if err := r.ParseForm(); err != nil {
			w.WriteHeader(200)
			// components.SignInForm(components.RenderAuthError(&[]utils.AuthError{{Err: err}})).Render(r.Context(), w)
			return
		}

		cred := auth.SignInCredentials{
			User: r.FormValue("user"),
			Pass: r.FormValue("password"),
		}

		user, errs := cred.SignIn()

		if errs != nil {
			w.WriteHeader(200)
			// components.SignInForm(components.RenderAuthError(errs)).Render(r.Context(), w)
			return
		}

		access, refresh, err := auth.NewTokenPair(
			&auth.JwtParams{
				Username: user.Username,
				UserId:   user.ID,
				Family:   uuid.New(),
			})

		if err != nil {
			w.WriteHeader(200)
			// components.SignInForm(components.RenderAuthError(errs)).Render(r.Context(), w)
			return
		}

		auth.SetTokenCookies(w, access, refresh)
		HandleRedirect(w, r, fmt.Sprintf("/profile/%s", claims.Username), nil)
		return
	}

	if r.Method == "GET" {
		w.WriteHeader(200)
		// components.SignIn().Render(r.Context(), w)
		return
	}
}

func SignOut(w http.ResponseWriter, r *http.Request) {
	auth.DeleteJwtCookies(w)
	HandleRedirect(w, r, "/", auth.ErrUserSignedOut)
}
