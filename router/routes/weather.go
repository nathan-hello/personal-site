package routes

import (
	"net/http"
	"os/exec"
	"strings"
)

func Weather(w http.ResponseWriter, r *http.Request) {
    if r.Method != "GET" {
        w.WriteHeader(http.StatusMethodNotAllowed)
        return
    }

    pathParts := strings.Split(r.URL.Path, "/")
    var location string
    if len(pathParts) > 2 {
        location = pathParts[2]
    } else {
        location = "London"
    }

    cmd := exec.Command("wego", location)
    
    output, err := cmd.Output()
    if err != nil {
        http.Error(w, "Error executing weather command: "+err.Error(), http.StatusInternalServerError)
        return
    }
    
    w.Header().Set("Content-Type", "text/plain")
    w.Write(output)
}
