package routes

import (
	"net/http"
	"os/exec"
)

func Weather(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	cmd := exec.Command("wego")
	
	output, err := cmd.Output()
	if err != nil {
		http.Error(w, "Error executing weather command: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "text/plain")
	
	w.Write(output)
}
