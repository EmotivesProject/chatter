package auth

import (
	"chatter/internal/db"
	"chatter/internal/logger"
	"chatter/model"
	"errors"
	"time"

	"github.com/gocql/gocql"
)

func CreateToken(username string) (model.Token, error) {
	dbSession := db.GetSession()
	token := model.Token{}
	iterable := dbSession.Query("SELECT * FROM users WHERE username = ? LIMIT 1", username).Consistency(gocql.One).Iter()
	defer iterable.Close()

	var user model.ShortenedUser
	found := false
	m := map[string]interface{}{}
	for iterable.MapScan(m) {
		found = true
		user.ID = m["id"].(gocql.UUID)
		user.Name = m["name"].(string)
		user.Username = m["username"].(string)
		user.Token = m["user_token"].(gocql.UUID)
	}

	if !found {
		return token, errors.New("User not found")
	}

	// Create and set items for return statement
	expiration := time.Now().Add(5 * time.Minute).Unix()
	token.Token = gocql.TimeUUID()
	token.Expiration = expiration
	token.Username = username

	err := dbSession.Query("UPDATE users SET user_token=? WHERE id=?;", token.Token, user.ID).Exec()
	if err != nil {
		return token, err
	}

	return token, nil
}

func ValidateToken(token string) (model.ShortenedUser, error) {
	dbSession := db.GetSession()
	logger.Infof("Trying to find token %s", token)
	iterable := dbSession.Query("SELECT * FROM users WHERE user_token = ? LIMIT 1", token).Consistency(gocql.One).Iter()
	defer iterable.Close()
	var user model.ShortenedUser

	found := false
	m := map[string]interface{}{}
	for iterable.MapScan(m) {
		found = true
		user.ID = m["id"].(gocql.UUID)
		user.Name = m["name"].(string)
		user.Username = m["username"].(string)
		user.Token = m["user_token"].(gocql.UUID)
	}

	if !found {
		return user, errors.New("User not found")
	}

	err := dbSession.Query("UPDATE users SET user_token = null WHERE id = ?", user.ID).Exec()
	return user, err
}
