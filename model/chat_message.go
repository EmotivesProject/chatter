package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ChatMessage struct {
	ID           primitive.ObjectID `bson:"_id" json:"id"`
	UsernameFrom string             `bson:"username_from" json:"username_from"`
	UsernameTo   string             `bson:"username_to" json:"username_to"`
	Message      string             `bson:"message,omitempty" json:"message,omitempty"`
	ImagePath    string             `bson:"image_path,omitempty" json:"image_path,omitempty"`
	Created      time.Time          `bson:"created" json:"created"`
}

func (c *ChatMessage) FillMessage() {
	c.ID = primitive.NewObjectID()
	c.Created = time.Now()
}

func (c ChatMessage) Validate() bool {
	if c.Message == "" && c.ImagePath == "" {
		return false
	}
	return true
}
