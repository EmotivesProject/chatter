package model

import (
	"time"
)

type ChatMessage struct {
	ID           int       `json:"id"`
	UsernameFrom string    `json:"username_from"`
	UsernameTo   string    `json:"username_to"`
	Message      string    `json:"message,omitempty"`
	ImagePath    string    `json:"image_path,omitempty"`
	Created      time.Time `json:"created"`
}

func (c *ChatMessage) FillMessage() {
	c.Created = time.Now()
}

func (c ChatMessage) Validate() bool {
	if c.Message == "" && c.ImagePath == "" {
		return false
	}

	return true
}
