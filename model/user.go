package model

type User struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	UserGroup string `json:"user_group"`
}
