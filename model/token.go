package model

type Token struct {
	Token      string `json:"token"`
	Expiration int64  `json:"expires"`
}
