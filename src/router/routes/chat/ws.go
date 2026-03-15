package chat

import (
	"bytes"
	"log"
	"net/http"
	"sync"

	gws "github.com/gorilla/websocket"
	"github.com/nathan-hello/personal-site/src/auth"
	"github.com/nathan-hello/personal-site/src/components"
	"github.com/nathan-hello/personal-site/src/db"
)

var upgrader = gws.Upgrader{CheckOrigin: func(r *http.Request) bool {
	return true
},
}

type Manager struct {
	clients map[*gws.Conn]bool
	lock    sync.Mutex
}

func (m *Manager) AddClient(c *gws.Conn) {

	m.lock.Lock()
	defer m.lock.Unlock()
	m.clients[c] = true
}

func (m *Manager) RemoveClient(c *gws.Conn) {
	m.lock.Lock()
	defer m.lock.Unlock()

	if _, ok := m.clients[c]; ok {
		delete(m.clients, c)
		c.Close()
	}
}

func (m *Manager) BroadcastMessage(message []byte) {
	m.lock.Lock()
	defer m.lock.Unlock()
	for c := range m.clients {
		if err := c.WriteMessage(gws.TextMessage, message); err != nil {
			log.Println(err)
			delete(m.clients, c)
			c.Close()
		}
	}
}

var manager = Manager{
	clients: make(map[*gws.Conn]bool),
}

func ChatSocket(w http.ResponseWriter, r *http.Request) {
	user, _ := auth.GetProfile(r)
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()
	manager.AddClient(conn)
	defer manager.RemoveClient(conn)

	for {
		_, clientMsg, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		msg, err := newChatFromBytes(clientMsg, user.Username, user.ID, user.Color)
		if err != nil {
			log.Println(err)
			w.Write([]byte(err.Error()))
		}

		db.Db().InsertMessage(
			r.Context(),
			db.InsertMessageParams{
				AuthorID:  user.ID,
				Message:   msg.Text,
				CreatedAt: msg.CreatedAt,
				RoomID:    DEFAULT_ROOM_ID,
			})

		buffMsg := &bytes.Buffer{}
		components.ChatMessage(msg).Render(r.Context(), buffMsg)
		manager.BroadcastMessage(buffMsg.Bytes())
	}
}
