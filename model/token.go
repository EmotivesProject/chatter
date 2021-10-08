package model

import "time"

type Token struct {
	Token      string    `json:"token"`
	Username   string    `json:"username"`
	Expiration time.Time `json:"expires"`
}
