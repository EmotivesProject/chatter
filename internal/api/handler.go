package api

import (
	"chatter/internal/auth"
	"chatter/internal/connections"
	"chatter/internal/db"
	"chatter/internal/messages"
	"chatter/model"
	"encoding/json"
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

func createUser(w http.ResponseWriter, r *http.Request) {
	user := model.User{}

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, false, http.StatusBadRequest, response.Message{Message: err.Error()})

		return
	}

	_, err = db.CreateUser(r.Context(), user)
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, false, http.StatusBadRequest, response.Message{Message: err.Error()})

		return
	}

	logger.Infof("Created user for chatter %s", user.Username)
	response.MessageResponseJSON(w, false, http.StatusCreated, response.Message{Message: "Created user"})
}

func createTocken(w http.ResponseWriter, r *http.Request) {
	user, err := getUsernameAndGroup(r)
	if err != nil {
		logger.Error(messages.ErrFailedToType)
		response.MessageResponseJSON(
			w,
			false,
			http.StatusInternalServerError,
			response.Message{Message: messages.ErrFailedToType.Error()},
		)

		return
	}

	token, err := auth.CreateToken(r.Context(), user, false)
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, false, http.StatusUnprocessableEntity, response.Message{Message: err.Error()})

		return
	}

	logger.Infof("Created token for user %s", user.Username)
	response.ResultResponseJSON(w, false, http.StatusCreated, token)
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	username, err := auth.ValidateToken(r.Context(), r.URL.Query().Get("token"))
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, false, http.StatusForbidden, response.Message{Message: err.Error()})

		return
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, false, http.StatusInternalServerError, response.Message{Message: err.Error()})

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
	user, err := getUsernameAndGroup(r)
	if err != nil {
		logger.Error(messages.ErrFailedToType)
		response.MessageResponseJSON(
			w,
			false,
			http.StatusInternalServerError,
			response.Message{Message: messages.ErrFailedToType.Error()},
		)

		return
	}

	offline := db.GetAllUsers(r.Context(), user.UserGroup)
	cons := connections.FilterOfflineUsers(*offline)

	logger.Infof("Sending user %s list of users", user.Username)

	response.ResultResponseJSON(w, false, http.StatusOK, cons)
}

func getMessages(w http.ResponseWriter, r *http.Request) {
	username := r.Context().Value(verification.UserID)

	from := r.URL.Query().Get("from")
	if username != from {
		response.MessageResponseJSON(
			w,
			false,
			http.StatusUnprocessableEntity,
			response.Message{Message: messages.WrongResponse},
		)

		return
	}

	to := r.URL.Query().Get("to")

	_, err := db.FindUserNoCreate(r.Context(), to)
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(
			w,
			false,
			http.StatusUnprocessableEntity,
			response.Message{Message: messages.WrongResponse},
		)

		return
	}

	skip := findSkip(r)
	messages := db.GetMessagesForUsers(r.Context(), from, to, skip)

	logger.Infof("Sending user %s messages talking to user %s", from, to)

	response.ResultResponseJSON(w, false, http.StatusOK, messages)
}
