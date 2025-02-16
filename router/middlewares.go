package router

import (
	"log"
	"net/http"
	"slices"
	"time"

	"github.com/justinas/alice"
)

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		defer func() {
			log.Printf("IP: %v, ROUTE REQUESTED: %v, RESPONSE TIME: %v\n", r.URL.User, r.URL.Path, time.Since(start))
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
