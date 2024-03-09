package routes

import "net/http"

func HandleRedirect(w http.ResponseWriter, r *http.Request, route string, err error) {
	http.Redirect(w, r, route, http.StatusSeeOther)
	if err != nil {
		w.WriteHeader(500)
	}

}
