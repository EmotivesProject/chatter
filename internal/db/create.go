package db

import (
	"chatter/model"

	"github.com/TomBowyerResearchProject/common/logger"

	"github.com/gocql/gocql"
)

func CreateUser(username string) error {
	session := GetSession()
	userID := gocql.TimeUUID()

	logger.Infof("Creating user %s", username)
	return session.Query("INSERT into users(id, username) values (?, ?);", userID, username).Exec()
}

func CreateMessage(msg model.ChatMessage) {
	session := GetSession()

	logger.Infof("Creating new message %s %s %s", msg.Message, msg.UsernameFrom, msg.UsernameTo)

	err := session.Query("INSERT into messages(id, created, message, username_from, username_to) values (?, ?, ?, ?, ?);", msg.ID, msg.Created, msg.Message, msg.UsernameFrom, msg.UsernameTo).Exec()
	if err != nil {
		logger.Error(err)
	}
}
