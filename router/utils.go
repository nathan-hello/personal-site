package router

import "net/http"

var StatusCodes = map[int]string{
	102: "https://http.cat/102.jpg",
	400: "https://http.cat/400.jpg",
	401: "https://http.cat/401.jpg",
	403: "https://http.cat/403.jpg",
	404: "https://http.cat/404.jpg",
	405: "https://http.cat/405.jpg",
	413: "https://http.cat/413.jpg",
	500: "https://http.cat/500.jpg",
}

func Index(b bool, dir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.URL.Path != "/" {
			http.Redirect(w, r, StatusCodes[404], http.StatusMovedPermanently)
			return
		}

		if b {
			fs := http.FileServer(http.Dir(dir))
			fs.ServeHTTP(w, r)
		}
	}
}

func ShowStatusCode(w http.ResponseWriter, r *http.Request, code int) {
	http.Redirect(w, r, StatusCodes[code], http.StatusMovedPermanently)
}
