package routes

import "net/http"

func ApiComments(w http.ResponseWriter, r *http.Request) {
        r.PathValue("id")
}
