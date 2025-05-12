package router

import (
	"context"
	"net/http"
	"slices"
	"time"

	"github.com/justinas/alice"
	"github.com/nathan-hello/personal-site/auth"
	"github.com/nathan-hello/personal-site/utils"
)

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		// TODO:
		// r.URL.Host is always blank
		// code is always blank even if i put in context (or did i????)
		// json is always blank even if i put in context (or did i????)

		defer func() {

			code, ok := r.Context().Value(utils.AnalyticsContextKeyHttpResponseStatus).(int)
			if !ok || code == 0 {
				code = 200
			}

			json, ok := r.Context().Value(utils.JsonContextKey).(string)
			if !ok || json == "" {
				json = "{}"
			}

			host := utils.RealIP(r)
			if host == "" {
				host = "0.0.0.0"
			}

			utils.HttpAnalytic(time.Now(), host, code, r.Method, r.URL.Path, start, json)
		}()
		next.ServeHTTP(w, r)
	})
}

func AllowMethods(methods ...string) alice.Constructor {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !slices.Contains(methods, r.Method) {
				w.WriteHeader(http.StatusMethodNotAllowed)
				w.Write([]byte("Method not allowed"))
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func CreateHeader(key string, value string) alice.Constructor {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add(key, value)
			next.ServeHTTP(w, r)
		})
	}
}

func RejectSubroute(path string) alice.Constructor {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != path {
				http.NotFound(w, r)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func InjectClaimsOnValidToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		access, ok := auth.ValidateJwtOrDelete(w, r)
		if !ok {
			next.ServeHTTP(w, r)
			return
		}

		user, _, err := auth.ParseToken(access)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		if user == nil {
			next.ServeHTTP(w, r)
			return
		}

		wrapReq := r.WithContext(context.WithValue(r.Context(), auth.UserContextKey, user))

		next.ServeHTTP(w, wrapReq)
	})
}

func ProtectedRoute(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// user, ok :=
		// if !ok || user != nil {
		// 	// TODO: redirect
		// }

	})
}
