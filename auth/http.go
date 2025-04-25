package auth

import (
	"net/http"
	"time"

	"github.com/nathan-hello/personal-site/db"
	"github.com/nathan-hello/personal-site/utils"
)

var UserContextKey = struct{}{}

func UserCtxDefaultAnon(r *http.Request) *db.SelectUserByIdRow {
   val := r.Context().Value(UserContextKey)
   if u, ok := val.(*db.SelectUserByIdRow); ok && u != nil {
       return u
   }
   return &db.SelectUserByIdRow{
       ID:              "anon",
       Email:           "anon",
       Username:        "Anonymous",
       GlobalChatColor: "purple-500",
   }
}

func SetTokenCookies(w http.ResponseWriter, a string, r string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    a,
		Expires:  time.Now().Add(utils.Env().REFRESH_EXPIRY_TIME),
		Secure:   true,
		HttpOnly: true,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    r,
		Expires:  time.Now().Add(utils.Env().REFRESH_EXPIRY_TIME),
		Secure:   true,
		HttpOnly: true,
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
	})
}

func getJwtsFromCookie(r *http.Request) (string, string, error) {
	access, err := r.Cookie("access_token")
	if err != nil {
		return "", "", err
	}

	refresh, err := r.Cookie("refresh_token")
	if err != nil {
		return "", "", err
	}

	return access.Value, refresh.Value, nil
}

func deleteCookie(w http.ResponseWriter, name string) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})
}

func DeleteJwtCookies(w http.ResponseWriter) {
	deleteCookie(w, "access_token")
	deleteCookie(w, "refresh_token")
}

func ValidateJwtOrDelete(w http.ResponseWriter, r *http.Request) (string, bool) {
	access, err := r.Cookie("access_token")
	if err != nil {
		if err == http.ErrNoCookie {
			return "", false
		}
		DeleteJwtCookies(w)
		return "", false
	}

	refresh, err := r.Cookie("refresh_token")
	if err != nil {
		if err == http.ErrNoCookie {
			return "", false
		}
		DeleteJwtCookies(w)
		return "", false
	}

	vAccess, vRefresh, err := validatePairOrRefresh(access.Value, refresh.Value)

	if err != nil {
		DeleteJwtCookies(w)
		return "", false
	}

	SetTokenCookies(w, vAccess, vRefresh)
	return vAccess, true
}
