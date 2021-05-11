package api

import (
	"chatter/internal/auth"
	"chatter/internal/connections"
	"chatter/internal/db"
	"chatter/internal/messages"
	"chatter/model"
	"net/http"

	"github.com/TomBowyerResearchProject/common/logger"
	"github.com/TomBowyerResearchProject/common/response"
	"github.com/TomBowyerResearchProject/common/verification"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func createTocken(w http.ResponseWriter, r *http.Request) {
	username, ok := r.Context().Value(verification.UserID).(string)
	if !ok {
		logger.Error(messages.ErrFailedToType)
		response.MessageResponseJSON(w, http.StatusBadRequest, response.Message{Message: messages.ErrFailedToType.Error()})

		return
	}

	token, err := auth.CreateToken(r.Context(), username, false)
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, http.StatusBadRequest, response.Message{Message: err.Error()})

		return
	}

	logger.Infof("Created token for user %s", username)
	response.ResultResponseJSON(w, http.StatusOK, token)
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	username, err := auth.ValidateToken(r.Context(), r.URL.Query().Get("token"))
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, http.StatusBadRequest, response.Message{Message: err.Error()})

		return
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, http.StatusBadRequest, response.Message{Message: err.Error()})

		return
	}
	// ensure connection close when function returns
	defer ws.Close()
	connections.Add(ws, username)
	logger.Infof("Connecting user %s", username)

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
	offline := db.GetAllUsers(r.Context())
	cons := connections.FilterOfflineUsers(*offline)
	response.ResultResponseJSON(w, http.StatusOK, cons)
}

func getMessages(w http.ResponseWriter, r *http.Request) {
	username := r.Context().Value(verification.UserID)

	from := r.URL.Query().Get("from")
	if username != from {
		response.MessageResponseJSON(w, http.StatusBadRequest, response.Message{Message: messages.WrongResponse})

		return
	}

	to := r.URL.Query().Get("to")
	skip := findSkip(r)
	messages := db.GetMessagesForUsers(r.Context(), from, to, skip)
	response.ResultResponseJSON(w, http.StatusOK, messages)
}
