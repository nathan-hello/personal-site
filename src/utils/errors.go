package utils

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

type AnalyticsContextKeyHttpResponseStatusType struct{}

var AnalyticsContextKeyHttpResponseStatus = AnalyticsContextKeyHttpResponseStatusType{}

type JsonContextKeyType struct{}

var JsonContextKey = JsonContextKeyType{}

var (
	ErrNoTextInChatMsg = errors.New("illegal message - no text in chat message")
)

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

func ShowStatusCode(w http.ResponseWriter, r *http.Request, code int) {
	ctx := context.WithValue(r.Context(), AnalyticsContextKeyHttpResponseStatus, StatusCodes[code])
	r = r.WithContext(ctx)

	http.Redirect(w, r, StatusCodes[code], http.StatusMovedPermanently)
}

func JsonAddError(r *http.Request, s string) *http.Request {
	ctx := context.WithValue(r.Context(), JsonContextKey, fmt.Sprintf("{\"error\":\"%s\"}", s))
	r = r.WithContext(ctx)
	return r
}
