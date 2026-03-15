package router

import (
	"net/http"
	"slices"
	"time"

	"github.com/justinas/alice"
	"github.com/nathan-hello/personal-site/src/auth"
	"github.com/nathan-hello/personal-site/src/utils"
)

func InjectUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		stack, cookie, ok := auth.GetSessionFromRequest(r)
		if !ok {
			next.ServeHTTP(w, r)
			return
		}

		if cookie != "" {
			w.Header().Set("Set-Cookie", cookie)
		}

		wrap := auth.InjectContext(r, stack)

		next.ServeHTTP(w, wrap)
	})
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
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
