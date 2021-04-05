package db

import (
	"chatter/model"
	"fmt"

	"github.com/gocql/gocql"
)

func CreateUser(user *model.ShortenedUser) error {
	session := GetSession()
	user.ID = gocql.TimeUUID()

	return session.Query("INSERT into users(id, name, username) values (?, ?, ?);", user.ID, user.Name, user.Username).Exec()
}

func CreateMessage(msg model.ChatMessage) {
	session := GetSession()

	err := session.Query("INSERT into messages(id, created, message, username_from, username_to) values (?, ?, ?, ?, ?);", msg.ID, msg.Created, msg.Message, msg.UsernameFrom, msg.UsernameTo).Exec()
	if err != nil {
		fmt.Println(err.Error())
	}
}
