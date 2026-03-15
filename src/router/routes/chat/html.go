package chat

import (
	"fmt"
	"log"
	"net/http"

	"github.com/nathan-hello/personal-site/src/auth"
	"github.com/nathan-hello/personal-site/src/components"
	"github.com/nathan-hello/personal-site/src/db"
	"github.com/nathan-hello/personal-site/src/utils"
)

func Chat(w http.ResponseWriter, r *http.Request) {
	user, _ := auth.GetProfile(r)
	embed := r.URL.Query().Get("embed") == "true"

	if r.Method == "GET" {
		fmt.Printf("got chat get\n")
		recents, err := db.GetMessagesByChatroom(
			r.Context(),
			db.SelectMessagesByChatroomParams{
				RoomID: DEFAULT_ROOM_ID,
				Limit:  10,
			})
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var renderedMessages []*utils.ChatMessage
		for _, msg := range recents {
			m := &utils.ChatMessage{
				UserId:    msg.Message.AuthorID,
				Text:      msg.Message.Message,
				Color:     msg.Profile.Color,
				CreatedAt: msg.Message.CreatedAt,
			}
			renderedMessages = append(renderedMessages, m)
		}

		components.ChatRoot(user, embed, renderedMessages).Render(r.Context(), w)
	}
}
