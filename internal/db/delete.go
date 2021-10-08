package db

import (
	"context"

	commonPostgres "github.com/TomBowyerResearchProject/common/postgres"
)

func DeleteToken(ctx context.Context, token string) error {
	connection := commonPostgres.GetDatabase()

	_, err := connection.Exec(
		ctx,
		`DELETE FROM tokens WHERE token = $1`,
		token,
	)

	return err
}
