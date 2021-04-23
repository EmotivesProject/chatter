package connections

import (
	"chatter/internal/db"
	"chatter/model"
	"io"
	"sync"

	"github.com/TomBowyerResearchProject/common/logger"

	"github.com/gorilla/websocket"
)

var (
	mapMutex    = sync.RWMutex{}
	connections = make(map[*websocket.Conn]string)
	clients     = make(map[string]*websocket.Conn)
	broadcaster = make(chan model.ChatMessage)
)

func Add(ws *websocket.Conn, username string) {
	logger.Infof("Adding %s", username)
	mapMutex.Lock()
	connections[ws] = username
	clients[username] = ws
	mapMutex.Unlock()
	notifyOfConnectionUpdate(username, true)
}

func Remove(ws *websocket.Conn) {
	mapMutex.Lock()
	username := connections[ws]
	logger.Infof("removing %s", username) // shouldn't be logging in a mutex lock
	delete(connections, ws)
	delete(clients, username)
	mapMutex.Unlock()
	notifyOfConnectionUpdate(username, false)
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

func notifyOfConnectionUpdate(username string, active bool) {
	connection := model.Connection{
		Username: username,
		Active:   active,
	}

	for i, v := range clients {
		if i != username {
			messageClient(v, connection)
		}
	}
}

func HandleMessages() {
	for {
		msg := <-broadcaster
		msg.FillMessage()
		messageClients(msg)
	}
}

func messageClients(msg model.ChatMessage) {
	_, err := db.CreateMessage(msg)
	if err != nil {
		logger.Error(err)
	}

	to := clients[msg.UsernameTo]
	if to != nil {
		messageClient(to, msg)
	}

	from := clients[msg.UsernameFrom]
	if from != nil {
		messageClient(from, msg)
	}
}

func messageClient(client *websocket.Conn, msg interface{}) {
	err := client.WriteJSON(msg)
	if err != nil && unsafeError(err) {
		logger.Error(err)
		defer client.Close()
		Remove(client)
	}
}

func unsafeError(err error) bool {
	return !websocket.IsCloseError(err, websocket.CloseGoingAway) && err != io.EOF
}
