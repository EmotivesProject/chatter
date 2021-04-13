package api

import (
	"chatter/internal/auth"
	"chatter/internal/connections"
	"chatter/internal/db"
	"chatter/internal/logger"
	"chatter/internal/messages"
	"chatter/model"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func healthz(w http.ResponseWriter, r *http.Request) {
	messageResponseJSON(w, http.StatusOK, model.Message{Message: messages.MsgHealthOK})
}

func createTocken(w http.ResponseWriter, r *http.Request) {
	username := r.Context().Value(userID)
	token, err := auth.CreateToken(fmt.Sprintf("%v", username))
	if err != nil {
		logger.Error(err)
		messageResponseJSON(w, http.StatusBadRequest, model.Message{Message: err.Error()})
		return
	}
	logger.Infof("Created token for user %s", username)
	resultResponseJSON(w, http.StatusOK, token)
}

func createUser(w http.ResponseWriter, r *http.Request) {
	user := &model.ShortenedUser{}
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		logger.Error(err)
		messageResponseJSON(w, http.StatusBadRequest, model.Message{Message: err.Error()})
		return
	}

	err = db.CreateUser(user)
	if err != nil {
		logger.Error(err)
		messageResponseJSON(w, http.StatusBadRequest, model.Message{Message: err.Error()})
		return
	}

	logger.Info("Created new user")

	resultResponseJSON(w, http.StatusOK, user)
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	user, err := auth.ValidateToken(r.URL.Query().Get("token"))
	if err != nil {
		logger.Error(err)
		messageResponseJSON(w, http.StatusBadRequest, model.Message{Message: err.Error()})
		return
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Error(err)
		messageResponseJSON(w, http.StatusBadRequest, model.Message{Message: err.Error()})
		return
	}
	// ensure connection close when function returns
	defer ws.Close()
	connections.Add(ws, user.Username)
	logger.Infof("Connecting user %s", user.Username)

	for {
		var msg model.ChatMessage
		// Read in a new message as JSON and map it to a Message object
		err := ws.ReadJSON(&msg)
		if err != nil {
			connections.Remove(ws)
			break
		}
		// send new message to the channel
		connections.NewMessage(msg)
	}
}

func getConnectedUsers(w http.ResponseWriter, r *http.Request) {
	offline := db.GetAllUsers()
	cons := connections.FilterOfflineUsers(offline)
	resultResponseJSON(w, http.StatusOK, cons)
}

func getMessages(w http.ResponseWriter, r *http.Request) {
	username := r.Context().Value(userID)
	from := r.URL.Query().Get("from")
	if username != from {
		messageResponseJSON(w, http.StatusBadRequest, model.Message{Message: "Wrong"})
		return
	}
	to := r.URL.Query().Get("to")
	limit := findLimit(r)
	messages := db.GetMessagesForUsers(from, to, limit)
	resultResponseJSON(w, http.StatusOK, messages)
}
