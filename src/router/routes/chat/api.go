package chat

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/nathan-hello/personal-site/src/auth"
	"github.com/nathan-hello/personal-site/src/components"
	"github.com/nathan-hello/personal-site/src/db"
	"github.com/nathan-hello/personal-site/src/utils"
)

const DEFAULT_ROOM_ID = 1

func ApiChat(w http.ResponseWriter, r *http.Request) {
	user, _ := auth.GetProfile(r)
	htmlResponse := r.Header.Get("Content-Type") == "text/html"
	jsonResponse := r.Header.Get("Content-Type") == "application/json"
	if !htmlResponse && !jsonResponse {
		htmlResponse = true
	}

	if r.Method == "POST" {
		body, err := io.ReadAll(r.Body)

		if err != nil {
			fmt.Fprintf(w, "{error: \"%v\"}", err)
		}

		c, err := newChatFromBytes(body, user.Username, user.ID, user.Color)
		if err != nil {
			fmt.Fprintf(w, "{error: \"%v\"}\n", err)
			return
		}

		var resp []byte

		// We need to send html to subscribers no matter what
		htmlMsg := bytes.Buffer{}
		components.ChatMessage(c).Render(r.Context(), &htmlMsg)
		manager.BroadcastMessage(htmlMsg.Bytes())

		err = db.Db().InsertMessage(
			r.Context(),
			db.InsertMessageParams{
				AuthorID:  user.ID,
				Message:   c.Text,
				RoomID:    DEFAULT_ROOM_ID,
				CreatedAt: c.CreatedAt,
			})
		if err != nil {
			log.Println(err)
		}

		if htmlResponse {
			resp = htmlMsg.Bytes()
		}
		if jsonResponse {
			jason, err := json.Marshal(c)
			if err != nil {
				fmt.Fprintf(w, "{error: \"%v\"}", err)
			}
			resp = jason
		}
		w.Write(resp)
	}
}

type rawChatMessage struct {
	Text string `json:"msg-text"`
}

var ErrBadChatMsg = errors.New("not enough fields in chat message")

func newChatFromBytes(bits []byte, username string, userId string, color string) (*utils.ChatMessage, error) {
	var raw rawChatMessage
	err := json.Unmarshal(bits, &raw)
	if username == "" || userId == "" || color == "" {
		return nil, ErrBadChatMsg
	}
	if err != nil {
		return nil, err
	}
	if raw.Text == "" {
		return nil, utils.ErrNoTextInChatMsg
	}
	return &utils.ChatMessage{
		UserId:    userId,
		Username:  username,
		Text:      raw.Text,
		Color:     color,
		CreatedAt: time.Now().UTC(),
	}, nil
}
