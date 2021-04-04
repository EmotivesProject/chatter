package model

import (
	"time"

	"github.com/gocql/gocql"
)

type ChatMessage struct {
	ID           gocql.UUID `json:"id"`
	UsernameFrom string     `json:"username_from"`
	UsernameTo   string     `json:"username_to"`
	Message      string     `json:"message"`
	Created      time.Time  `json:"created"`
}

func (c *ChatMessage) FillMessage() *ChatMessage {
	c.ID = gocql.TimeUUID()
	c.Created = time.Now()
	return c
}
