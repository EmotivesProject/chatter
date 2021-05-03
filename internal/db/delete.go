package db

import (
	"context"

	commonMongo "github.com/TomBowyerResearchProject/common/mongo"
	"gopkg.in/mgo.v2/bson"
)

func DeleteToken(token string) error {
	db := commonMongo.GetDatabase()
	tokenCollection := db.Collection(TokensCollection)
	_, err := tokenCollection.DeleteOne(context.TODO(), bson.M{"token": token})

	return err
}
