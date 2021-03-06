package connections

import (
	"chatter/internal/db"
	"chatter/internal/send"
	"chatter/model"
	"context"
	"errors"
	"io"
	"sync"

	"github.com/EmotivesProject/common/logger"

	"github.com/gorilla/websocket"
)

var (
	mapMutex    = sync.RWMutex{}
	connections = make(map[*websocket.Conn]string)
	clients     = make(map[string]*websocket.Conn)
	broadcaster = make(chan model.ChatMessage)
)

func Add(ws *websocket.Conn, username string) {
	logger.Infof("Adding %s to websocket connections", username)
	mapMutex.Lock()
	connections[ws] = username
	clients[username] = ws
	mapMutex.Unlock()
	notifyOfConnectionUpdate(username, true)
}

func Remove(ws *websocket.Conn) {
	mapMutex.Lock()
	username := connections[ws]
	delete(connections, ws)
	delete(clients, username)
	mapMutex.Unlock()
	logger.Infof("removed %s from websocket connections", username)
	notifyOfConnectionUpdate(username, false)
}

func NewMessage(message model.ChatMessage) {
	broadcaster <- message
}

func FilterOfflineUsers(offlineUsers []model.Connection) []model.Connection {
	mapMutex.Lock()
	for i, v := range offlineUsers {
		if clients[v.Username] != nil {
			offlineUsers[i].Active = true
		}
	}
	mapMutex.Unlock()

	return offlineUsers
}

func notifyOfConnectionUpdate(username string, active bool) {
	connection := model.Connection{
		Username: username,
		Active:   active,
	}

	for i, v := range clients {
		if i != username && db.IsSameGroup(username, i) {
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
	logger.Infof("Handling message %s from %s to %s", msg.Message, msg.UsernameFrom, msg.UsernameTo)

	send.MessageNotification(msg.UsernameFrom, msg.UsernameTo, msg.Message)

	if _, err := db.CreateMessage(context.Background(), msg); err != nil {
		logger.Error(err)
	}

	if to := clients[msg.UsernameTo]; to != nil {
		messageClient(to, msg)
	}

	from := clients[msg.UsernameFrom]
	if from != nil {
		messageClient(from, msg)
	}
}

func messageClient(client *websocket.Conn, msg interface{}) {
	if err := client.WriteJSON(msg); err != nil && unsafeError(err) {
		logger.Error(err)

		defer client.Close()

		Remove(client)
	}
}

func unsafeError(err error) bool {
	return !websocket.IsCloseError(err, websocket.CloseGoingAway) && !errors.Is(err, io.EOF)
}
