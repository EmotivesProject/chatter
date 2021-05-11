package db

import (
	"context"

	commonMongo "github.com/TomBowyerResearchProject/common/mongo"
	"gopkg.in/mgo.v2/bson"
)

func DeleteToken(ctx context.Context, token string) error {
	db := commonMongo.GetDatabase()
	tokenCollection := db.Collection(TokensCollection)
	_, err := tokenCollection.DeleteOne(ctx, bson.M{"token": token})

	return err
}
