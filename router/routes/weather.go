package routes

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strings"

	"github.com/nathan-hello/personal-site/utils"
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
		addr := utils.RealIP(r)
	    location = getLocation(addr)
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

type ipAPIResponse struct {
    Status     string `json:"status"`
    City       string `json:"city"`
    RegionName string `json:"regionName"`
    Country    string `json:"country"`
}

// getLocation calls ip-api.com and returns "City, Region, Country"
func getLocation(ip string) string {
    url := fmt.Sprintf(
        "http://ip-api.com/json/%s?fields=status,city,regionName,country",
        ip,
    )
    resp, err := http.Get(url)
    if err != nil {
        log.Printf("geolookup HTTP error: %v", err)
        return ip
    }
    defer resp.Body.Close()

    var data ipAPIResponse
    if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
        log.Printf("geolookup JSON error: %v", err)
        return ip
    }
    if data.Status != "success" {
        log.Printf("geolookup failed for %s: %+v", ip, data)
        return ip
    }

    parts := []string{}
    if data.City != "" {
        parts = append(parts, data.City)
    }
    if data.RegionName != "" {
        parts = append(parts, data.RegionName)
    }
    if data.Country != "" {
        parts = append(parts, data.Country)
    }
    return strings.Join(parts, ", ")
}
