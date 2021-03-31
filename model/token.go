package model

import "github.com/gocql/gocql"

type Token struct {
	Token      gocql.UUID `json:"token"`
	Expiration int64      `json:"expires"`
}
