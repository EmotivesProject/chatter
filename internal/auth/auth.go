package auth

import (
	"chatter/internal/db"
	"chatter/internal/messages"
	"chatter/model"
	"context"
	"crypto/rand"
	"encoding/base64"
	"time"
)

const (
	expiration  = 5
	tokenLength = 20
)

func CreateToken(ctx context.Context, username string, previouslyCalled bool) (model.Token, error) {
	var token model.Token

	_, err := db.FindUser(ctx, username)
	if err != nil {
		return token, messages.ErrFailedUsername
	}

	generatedToken, err := generateRandomString(tokenLength)
	if err != nil {
		return token, err
	}

	// Create and set items for return statement
	expiration := time.Now().Add(expiration * time.Minute).Unix()
	token.Token = generatedToken
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

func generateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)

	_, err := rand.Read(b)

	return b, err
}

func generateRandomString(s int) (string, error) {
	b, err := generateRandomBytes(s)

	return base64.URLEncoding.EncodeToString(b), err
}
