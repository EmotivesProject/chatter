package api

import (
	"chatter/internal/auth"
	"chatter/internal/db"
	"chatter/internal/messages"
	"chatter/model"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"

	"github.com/gocql/gocql"
	"github.com/gorilla/websocket"
)

type ChatMessage struct {
	Username string `json:"username"`
	Text     string `json:"text"`
}

var (
	mapMuxtex   = sync.RWMutex{}
	connections = make(map[*websocket.Conn]string)
	clients     = make(map[string]*websocket.Conn)
	broadcaster = make(chan ChatMessage)
	upgrader    = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	session = db.Init()
)

func healthz(w http.ResponseWriter, r *http.Request) {
	messageResponseJSON(w, http.StatusOK, model.Message{Message: messages.MsgHealthOK})
}

func createTocken(w http.ResponseWriter, r *http.Request) {
	username := r.Context().Value(userID)
	token, err := auth.CreateToken(fmt.Sprintf("%v", username), session)
	if err != nil {
		messageResponseJSON(w, http.StatusBadRequest, model.Message{Message: err.Error()})
		return
	}
	resultResponseJSON(w, http.StatusOK, token)
}

func createUser(w http.ResponseWriter, r *http.Request) {
	user := &model.ShortenedUser{}
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		messageResponseJSON(w, http.StatusBadRequest, model.Message{Message: err.Error()})
		return
	}

	user.ID = gocql.TimeUUID()

	err = session.Query("INSERT into users(id, name, username) values (?, ?, ?);", user.ID, user.Name, user.Username).Exec()
	if err != nil {
		fmt.Println("Failed creating new user")
		messageResponseJSON(w, http.StatusBadRequest, model.Message{Message: err.Error()})
		return
	}

	fmt.Println("Created new user")

	resultResponseJSON(w, http.StatusOK, user)
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	user, err := auth.ValidateToken(r.URL.Query().Get("token"), session)
	if err != nil {
		messageResponseJSON(w, http.StatusBadRequest, model.Message{Message: err.Error()})
		return
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		messageResponseJSON(w, http.StatusBadRequest, model.Message{Message: err.Error()})
		return
	}
	// ensure connection close when function returns
	defer ws.Close()
	mapMuxtex.Lock()
	connections[ws] = user.Username
	clients[user.Username] = ws
	mapMuxtex.Unlock()
	fmt.Println("Connected new client")

	for {
		var msg ChatMessage
		// Read in a new message as JSON and map it to a Message object
		err := ws.ReadJSON(&msg)
		if err != nil {
			mapMuxtex.Lock()
			delete(connections, ws)
			delete(clients, user.Username)
			mapMuxtex.Unlock()
			break
		}
		// send new message to the channel
		broadcaster <- msg
	}
}

// If a message is sent while a client is closing, ignore the error
func unsafeError(err error) bool {
	return !websocket.IsCloseError(err, websocket.CloseGoingAway) && err != io.EOF
}

func HandleMessages() {
	for {
		// grab any next message from channel
		msg := <-broadcaster

		messageClients(msg)
	}
}

func messageClients(msg ChatMessage) {
	// send to every client currently connected
	mapMuxtex.Lock()
	for connection := range connections {
		messageClient(connection, msg)
	}
	mapMuxtex.Unlock()
}

func messageClient(client *websocket.Conn, msg ChatMessage) {
	err := client.WriteJSON(msg)
	if err != nil && unsafeError(err) {
		log.Printf("error: %v", err)
		client.Close()

		user := connections[client]
		delete(connections, client)
		delete(clients, user)
	}
}

func getConnectedUsers(w http.ResponseWriter, r *http.Request) {
	var number = 0
	mapMuxtex.Lock()
	var cons []model.Connection
	for _, v := range connections {
		number++
		connection := model.Connection{
			Username: v,
			Active:   true,
		}
		cons = append(cons, connection)
	}
	mapMuxtex.Unlock()
	fmt.Printf("SENDING OFF %d\n", number)
	resultResponseJSON(w, http.StatusOK, cons)
}
