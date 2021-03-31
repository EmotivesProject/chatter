package auth

import (
	"chatter/model"

	"github.com/google/uuid"
)

func CreateToken() model.Token {
	//expiration := time.Now().Add(5 * time.Minute).Unix()
	return model.Token{
		Token:      uuid.NewString(),
		Expiration: 1,
	}
}

func ValidateToken(token string) bool {
	return true
}
