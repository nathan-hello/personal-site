package router

import (
	"net/http"

	"github.com/nathan-hello/personal-site/utils"
)

func Index(b bool, dir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.URL.Path != "/" {
			http.Redirect(w, r, utils.StatusCodes[404], http.StatusMovedPermanently)
			return
		}

		if b {
			fs := http.FileServer(http.Dir(dir))
			fs.ServeHTTP(w, r)
		}
	}
}
