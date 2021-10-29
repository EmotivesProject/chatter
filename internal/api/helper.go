package api

import (
	"chatter/internal/messages"
	"chatter/model"
	"net/http"
	"strconv"

	"github.com/EmotivesProject/common/verification"
)

const (
	bit = 64
	ten = 10
)

func findSkip(r *http.Request) int64 {
	skipParam := r.URL.Query().Get("skip")
	if skipParam == "" {
		skipParam = "0"
	}

	limit, err := strconv.ParseInt(skipParam, ten, bit)
	if err != nil {
		return 0
	}

	return limit
}

func getUsernameAndGroup(r *http.Request) (model.User, error) {
	user := model.User{}

	username, ok := r.Context().Value(verification.UserID).(string)
	if !ok {
		return user, messages.ErrFailedToType
	}

	userGroup, ok := r.Context().Value(verification.UserGroup).(string)
	if !ok {
		return user, messages.ErrFailedToType
	}

	user.Username = username
	user.UserGroup = userGroup

	return user, nil
}
