package api

import (
	"chatter/model"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type ChatMessage struct {
	Username string `json:"username"`
	Text     string `json:"text"`
}

var (
	msgHealthOK = "Health ok"
	clients     = make(map[*websocket.Conn]bool)
	broadcaster = make(chan ChatMessage)
	upgrader    = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func healthz(w http.ResponseWriter, r *http.Request) {
	messageResponseJSON(w, http.StatusOK, model.Message{Message: msgHealthOK})
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	// ensure connection close when function returns
	defer ws.Close()
	clients[ws] = true

	for {
		var msg ChatMessage
		// Read in a new message as JSON and map it to a Message object
		err := ws.ReadJSON(&msg)
		if err != nil {
			delete(clients, ws)
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
	for client := range clients {
		messageClient(client, msg)
	}
}

func messageClient(client *websocket.Conn, msg ChatMessage) {
	err := client.WriteJSON(msg)
	if err != nil && unsafeError(err) {
		log.Printf("error: %v", err)
		client.Close()
		delete(clients, client)
	}
}
