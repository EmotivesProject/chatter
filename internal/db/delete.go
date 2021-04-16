package db

import (
	"context"

	"gopkg.in/mgo.v2/bson"
)

func DeleteToken(token string) error {
	db := GetDatabase()
	tokenCollection := db.Collection(TokensCollection)
	_, err := tokenCollection.DeleteOne(context.TODO(), bson.M{"token": token})
	return err
}
