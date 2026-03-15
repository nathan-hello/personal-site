package auth

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"math/rand"
	"net/http"
	"time"

	"github.com/nathan-hello/personal-site/src/db"
)

const COOKIE_PREFIX = "reluekiss"

type Session struct {
	ExpiresAt time.Time `json:"expiresAt"`
	Token     string    `json:"token"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	IPAddress string    `json:"ipAddress"`
	UserAgent string    `json:"userAgent"`
	UserID    string    `json:"userId"`
	ID        string    `json:"id"`
}

type User struct {
	Name             string    `json:"name"`
	Email            string    `json:"email"`
	EmailVerified    bool      `json:"emailVerified"`
	Image            *string   `json:"image"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
	Username         string    `json:"username"`
	DisplayUsername  string    `json:"displayUsername"`
	TwoFactorEnabled bool      `json:"twoFactorEnabled"`
	ID               string    `json:"id"`
}

type BAuthSession struct {
	Session Session `json:"session"`
	User    User    `json:"user"`
}

type tAuthKey uint64

var vAuthKey tAuthKey

func GetProfile(r *http.Request) (db.Profile, bool) {
	user, ok := GetUser(r)
	if !ok {
		return anonProfile(), false
	}

	profile, err := db.Conn.SelectUserProfileById(r.Context(), user.ID)
	if err == nil {
		return profile, true
	}
	if err == sql.ErrNoRows {
		profile, err = db.Conn.InsertUserProfile(r.Context(), db.InsertUserProfileParams{
			ID:       user.ID,
			Username: user.Username,
			Color:    color_map["green"],
		})
		if err != nil {
			return db.Profile{}, false
		}
	}

	return profile, true
}

func GetUser(r *http.Request) (User, bool) {
	ba, ok := GetSessionAndUser(r)
	if !ok {
		return User{}, false
	}
	return ba.User, true
}

func GetSession(r *http.Request) (Session, bool) {
	ba, ok := GetSessionAndUser(r)
	if !ok {
		return Session{}, false
	}
	return ba.Session, true
}

func GetSessionAndUser(r *http.Request) (BAuthSession, bool) {
	session, ok := r.Context().Value(vAuthKey).(*BAuthSession)
	if !ok {
		return BAuthSession{}, false
	}
	return *session, true
}

func InjectContext(r *http.Request, data BAuthSession) *http.Request {
	ctx := context.Background()
	ctx = context.WithValue(ctx, vAuthKey, data)

	wrap := r.WithContext(ctx)
	return wrap
}

func GetSessionFromRequest(r *http.Request) (BAuthSession, string, bool) {

	client := &http.Client{}
	var cookie *http.Cookie
	cookie, err := r.Cookie("__Secure." + COOKIE_PREFIX + ".session_token")
	if err != nil {
		cookie, err = r.Cookie(COOKIE_PREFIX + ".session_token")
		if err != nil {
			return BAuthSession{}, "", false
		}
	}

	req, err := http.NewRequest("GET", "https://reluekiss.com/api/auth/get-session", nil)
	if err != nil {
		return BAuthSession{}, "", false
	}

	req.Header.Set("Cookie: "+cookie.Name, cookie.Value)

	response, err := client.Do(req)
	if err != nil {
		return BAuthSession{}, "", false
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return BAuthSession{}, "", false
	}

	if bytes.Equal(body, []byte("null")) {
		return BAuthSession{}, "", false
	}

	var stack = BAuthSession{}

	err = json.Unmarshal(body, &stack)
	if err != nil {
		return BAuthSession{}, "", false
	}

	c := response.Header.Get("Set-Cookie")

	return stack, c, true
}

func anonProfile() db.Profile {

	colors := make([]string, 0, len(color_map))
	for _, v := range color_map {
		colors = append(colors, v)
	}

	num := rand.Intn(6)
	c := colors[num]

	return db.Profile{
		ID:       "anon",
		Username: "anon",
		Color:    c,
	}
}

var color_map map[string]string = map[string]string{
	"orange": "#f97316",
	"red":    "#ef4444",
	"yellow": "#eab308",
	"blue":   "#3b82f6",
	"pink":   "#ec4899",
	"green":  "#22c55e",
}
