package api

import (
	"chatter/model"
	"net/http"
)

var (
	msgHealthOK = "Health ok"
)

func healthz(w http.ResponseWriter, r *http.Request) {
	messageResponseJSON(w, http.StatusOK, model.Message{Message: msgHealthOK})
}
