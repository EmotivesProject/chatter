package db

import (
	"chatter/model"
	"time"

	"github.com/gocql/gocql"
)

const (
	PageLimit = 20
)

func GetAllUsers() []model.Connection {
	var userList []model.Connection
	m := map[string]interface{}{}
	session := GetSession()
	iterable := session.Query("SELECT username from users;").Iter()
	for iterable.MapScan(m) {
		userList = append(userList, model.Connection{
			Username: m["username"].(string),
			Active:   false,
		})
		m = map[string]interface{}{}
	}
	return userList
}

func GetMessagesForUsers(from, to string, begin int64) []model.ChatMessage {
	var chatList []model.ChatMessage
	m := map[string]interface{}{}
	session := GetSession()
	iterable := session.Query("select * from messages where username_to IN (?, ?) AND username_from IN (?, ?);", from, to, from, to).Iter()
	for iterable.MapScan(m) {
		chatList = append(chatList, model.ChatMessage{
			ID:           m["id"].(gocql.UUID),
			UsernameFrom: m["username_from"].(string),
			UsernameTo:   m["username_to"].(string),
			Message:      m["message"].(string),
			Created:      m["created"].(time.Time),
		})
		m = map[string]interface{}{}
	}
	return chatList
}
