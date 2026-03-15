package routes

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/a-h/templ"
	"github.com/nathan-hello/personal-site/src/components"
	"github.com/nathan-hello/personal-site/src/components/layouts"
	"github.com/nathan-hello/personal-site/src/utils"
)

func TvRoute(w http.ResponseWriter, r *http.Request) {
	isStreaming := isMediaMtxStreaming("mystream")

	child := templ.WithChildren(context.Background(), components.TvPage(isStreaming))

	layouts.BaseLayout(components.Header(utils.AsciiTv), components.Meta("NatTV", "Live streaming", "https://reluekiss.com/tv.html", nil)).Render(child, w)
}

func isMediaMtxStreaming(streamName string) bool {
	resp, err := http.Get("http://localhost:9997/v3/hlsmuxers/get/" + streamName)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false
	}

	// Try to unmarshal as success response
	var successResp struct {
		Path        string `json:"path"`
		Created     string `json:"created"`
		LastRequest string `json:"lastRequest"`
		BytesSent   int64  `json:"bytesSent"`
	}
	if err := json.Unmarshal(body, &successResp); err == nil {
		// Successfully unmarshaled as success -> stream exists
		return true
	}

	// Try to unmarshal as error response
	var errorResp struct {
		Status string `json:"status"`
		Error  string `json:"error"`
	}
	if err := json.Unmarshal(body, &errorResp); err == nil {
		if errorResp.Error == "muxer not found" {
			return false
		}
		// Other error -> treat as not streaming
		return false
	}

	// Could not unmarshal into either -> treat as not streaming
	return false
}
