package db

import (
	"chatter/model"
	"context"

	commonMongo "github.com/TomBowyerResearchProject/common/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateUser(username string) (*model.User, error) {
	user := model.User{
		Username: username,
	}
	user.ID = primitive.NewObjectID()

	_, err := insetIntoCollection(UsersCollection, user)

	return &user, err
}

func CreateMessage(msg model.ChatMessage) (*model.ChatMessage, error) {
	_, err := insetIntoCollection(MessageCollection, msg)

	return &msg, err
}

func CreateToken(token model.Token) (*model.Token, error) {
	_, err := insetIntoCollection(TokensCollection, token)

	return &token, err
}

func insetIntoCollection(collectionName string, document interface{}) (*mongo.InsertOneResult, error) {
	db := commonMongo.GetDatabase()
	collection := db.Collection(collectionName)

	return collection.InsertOne(context.TODO(), document)
}
