package db

import (
	"chatter/model"
	"context"

	commonMongo "github.com/TomBowyerResearchProject/common/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateUser(ctx context.Context, username string) (*model.User, error) {
	user := model.User{
		Username: username,
	}
	user.ID = primitive.NewObjectID()

	_, err := insetIntoCollection(ctx, UsersCollection, user)

	return &user, err
}

func CreateMessage(ctx context.Context, msg model.ChatMessage) (*model.ChatMessage, error) {
	_, err := insetIntoCollection(ctx, MessageCollection, msg)

	return &msg, err
}

func CreateToken(ctx context.Context, token model.Token) (*model.Token, error) {
	_, err := insetIntoCollection(ctx, TokensCollection, token)

	return &token, err
}

func insetIntoCollection(
	ctx context.Context,
	collectionName string,
	document interface{},
) (*mongo.InsertOneResult, error) {
	db := commonMongo.GetDatabase()
	collection := db.Collection(collectionName)

	return collection.InsertOne(ctx, document)
}
