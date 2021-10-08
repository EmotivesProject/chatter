package db

import (
	"chatter/model"
	"context"

	commonPostgres "github.com/TomBowyerResearchProject/common/postgres"
)

func CreateUser(ctx context.Context, username string) (*model.User, error) {
	user := model.User{
		Username: username,
	}

	connection := commonPostgres.GetDatabase()

	_, err := connection.Exec(
		ctx,
		`INSERT INTO users(username) VALUES ($1) RETURNING id`,
		user.Username,
	)

	return &user, err
}

func CreateMessage(ctx context.Context, msg model.ChatMessage) (*model.ChatMessage, error) {
	connection := commonPostgres.GetDatabase()

	err := connection.QueryRow(
		ctx,
		`INSERT INTO messages(username_from,username_to,message,image_path,created_at)
		VALUES ($1,$2,$3,$4,$5) RETURNING id`,
		msg.UsernameFrom,
		msg.UsernameTo,
		msg.Message,
		msg.ImagePath,
		msg.Created,
	).Scan(
		&msg.ID,
	)

	return &msg, err
}

func CreateToken(ctx context.Context, token model.Token) error {
	connection := commonPostgres.GetDatabase()

	_, err := connection.Exec(
		ctx,
		`INSERT INTO tokens(token,username,expiration)
		VALUES ($1,$2,$3)`,
		token.Token,
		token.Username,
		token.Expiration,
	)

	return err
}
