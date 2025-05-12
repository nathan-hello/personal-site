package routes

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"regexp"
	"strings"

	"github.com/nathan-hello/personal-site/utils"
)

func Weather(w http.ResponseWriter, r *http.Request) {
  if r.Method != http.MethodGet {
    w.WriteHeader(http.StatusMethodNotAllowed)
    return
  }

  pathParts := strings.Split(r.URL.Path, "/")
  var location string
  if len(pathParts) > 2 && pathParts[2] != "" {
    location = pathParts[2]
  } else {
    addr := utils.RealIP(r)
    location = getLocation(addr)
  }

  cmd := exec.Command("wego", location)
  rawOutput, err := cmd.Output()
  if err != nil {
    http.Error(w, "Error executing weather command: "+err.Error(),
      http.StatusInternalServerError)
    return
  }

  output := string(rawOutput)
  ua := r.Header.Get("User-Agent")
  if isBrowserUA(ua) {
    output = ansiRegexp.ReplaceAllString(output, "")
  }

  w.Header().Set("Content-Type", "text/plain; charset=utf-8")
  w.WriteHeader(http.StatusOK)
  w.Write([]byte(output))
}

var ansiRegexp = regexp.MustCompile("\x1b\\[[0-9;]*[a-zA-Z]")

func isBrowserUA(ua string) bool {
  ua = strings.ToLower(ua)
  return strings.Contains(ua, "mozilla") ||
    strings.Contains(ua, "chrome") ||
    strings.Contains(ua, "safari") ||
    strings.Contains(ua, "firefox") ||
    strings.Contains(ua, "edge") ||
    strings.Contains(ua, "opr") // Opera
}

type ipAPIResponse struct {
    Status     string `json:"status"`
    City       string `json:"city"`
}

func getLocation(ip string) string {
    url := fmt.Sprintf(
        "http://ip-api.com/json/%s?fields=status,city",
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

    return data.City
}
