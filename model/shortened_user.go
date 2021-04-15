package model

import "github.com/gocql/gocql"

//User struct declaration
type ShortenedUser struct {
	ID       gocql.UUID `json:"id"`
	Username string     `json:"username"`
	Token    gocql.UUID `json:"token"`
}
