package model

type Token struct {
	Token      string `bson:"token" json:"token"`
	Username   string `bson:"username" json:"username"`
	Expiration int64  `bson:"expires" json:"expires"`
}
