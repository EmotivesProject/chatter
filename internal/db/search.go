package db

import (
	"chatter/model"
	"context"
	"errors"

	"github.com/TomBowyerResearchProject/common/logger"
	commonPostgres "github.com/TomBowyerResearchProject/common/postgres"
	"github.com/jackc/pgx/v4"
)

func FindUser(ctx context.Context, username string) (*model.User, error) {
	user := model.User{}

	connection := commonPostgres.GetDatabase()

	err := connection.QueryRow(
		ctx,
		`SELECT * FROM users
		WHERE username = $1`,
		username,
	).Scan(
		&user.ID,
		&user.Username,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		userdef, err := CreateUser(ctx, username)

		return userdef, err
	} else if err != nil {
		logger.Error(err)
	}

	return &user, err
}

func FindUserNoCreate(ctx context.Context, username string) (*model.User, error) {
	user := model.User{}

	connection := commonPostgres.GetDatabase()

	err := connection.QueryRow(
		ctx,
		`SELECT * FROM users
		WHERE username = $1`,
		username,
	).Scan(
		&user.ID,
		&user.Username,
	)

	return &user, err
}

func FindToken(ctx context.Context, token string) (*model.Token, error) {
	tokenObj := model.Token{}

	connection := commonPostgres.GetDatabase()

	err := connection.QueryRow(
		ctx,
		`SELECT * FROM tokens
		WHERE token = $1`,
		token,
	).Scan(
		&tokenObj.Token,
		&tokenObj.Username,
		&tokenObj.Expiration,
	)

	return &tokenObj, err
}

func GetAllUsers(ctx context.Context) *[]model.Connection {
	userList := make([]model.Connection, 0)

	connection := commonPostgres.GetDatabase()

	rows, err := connection.Query(
		ctx,
		`SELECT username FROM users`,
	)
	if err != nil {
		return &userList
	}

	for rows.Next() {
		var connection model.Connection

		err := rows.Scan(
			&connection.Username,
		)
		if err != nil {
			logger.Error(err)
			continue
		}

		connection.Active = false

		userList = append(userList, connection)
	}

	return &userList
}

func GetMessagesForUsers(ctx context.Context, from, to string, skip int64) *[]model.ChatMessage {
	chatList := make([]model.ChatMessage, 0)

	connection := commonPostgres.GetDatabase()

	rows, err := connection.Query(
		ctx,
		`SELECT * FROM messages
		WHERE username_from = $1 AND username_to = $2
		OR username_to = $1 AND username_from = $2`,
		from,
		to,
	)
	if err != nil {
		return &chatList
	}

	for rows.Next() {
		var chatmessage model.ChatMessage

		err := rows.Scan(
			&chatmessage.ID,
			&chatmessage.UsernameFrom,
			&chatmessage.UsernameTo,
			&chatmessage.Message,
			&chatmessage.ImagePath,
			&chatmessage.Created,
		)
		if err != nil {
			logger.Error(err)
			continue
		}

		chatList = append(chatList, chatmessage)
	}

	return &chatList
}
