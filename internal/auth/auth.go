package auth

import (
	"chatter/internal/db"
	"chatter/internal/messages"
	"chatter/model"
	"context"
	"math/rand"
	"time"
)

const (
	letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	expiration  = 5
	tokenLength = 20
)

func CreateToken(ctx context.Context, username string, previouslyCalled bool) (model.Token, error) {
	var token model.Token

	_, err := db.FindUser(ctx, username)
	if err != nil {
		return token, messages.ErrFailedUsername
	}

	// Create and set items for return statement
	expiration := time.Now().Add(expiration * time.Minute).Unix()
	token.Token = RandStringBytes(tokenLength)
	token.Expiration = expiration
	token.Username = username

	_, err = db.CreateToken(ctx, token)
	if err != nil {
		return token, err
	}

	return token, nil
}

func ValidateToken(ctx context.Context, token string) (string, error) {
	tokenObj, err := db.FindToken(ctx, token)
	if err != nil {
		return "", err
	}

	err = db.DeleteToken(ctx, tokenObj.Token)

	return tokenObj.Username, err
}

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		// nolint: gosec
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}

	return string(b)
}
