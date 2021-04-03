package model

import "github.com/gocql/gocql"

type Token struct {
	Token      gocql.UUID `json:"token"`
	Username   string     `json:"username"`
	Expiration int64      `json:"expires"`
}
