package connections

import (
	"chatter/internal/db"
	"chatter/model"
	"io"
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	mapMutex    = sync.RWMutex{}
	connections = make(map[*websocket.Conn]string)
	clients     = make(map[string]*websocket.Conn)
	broadcaster = make(chan model.ChatMessage)
)

func Add(ws *websocket.Conn, username string) {
	mapMutex.Lock()
	connections[ws] = username
	clients[username] = ws
	mapMutex.Unlock()
}

func Remove(ws *websocket.Conn) {
	mapMutex.Lock()
	user := connections[ws]
	delete(connections, ws)
	delete(clients, user)
	mapMutex.Unlock()
}

func NewMessage(message model.ChatMessage) {
	broadcaster <- message
}

func FilterOfflineUsers(OfflineUsers []model.Connection) []model.Connection {
	mapMutex.Lock()
	for i, v := range OfflineUsers {
		if clients[v.Username] != nil {
			OfflineUsers[i].Active = true
		}
	}
	mapMutex.Unlock()
	return OfflineUsers
}

func HandleMessages() {
	for {
		msg := <-broadcaster
		msg = *msg.FillMessage()
		messageClients(msg)
	}
}

func messageClients(msg model.ChatMessage) {
	db.CreateMessage(msg)
	to := clients[msg.UsernameTo]
	if to != nil {
		messageClient(to, msg)
	}
	from := clients[msg.UsernameFrom]
	if from != nil {
		messageClient(from, msg)
	}
}

func messageClient(client *websocket.Conn, msg model.ChatMessage) {
	err := client.WriteJSON(msg)
	if err != nil && unsafeError(err) {
		log.Printf("error: %v", err)
		defer client.Close()
		Remove(client)
	}
}

func unsafeError(err error) bool {
	return !websocket.IsCloseError(err, websocket.CloseGoingAway) && err != io.EOF
}
