package db

import (
	"chatter/internal/messages"
	"chatter/model"
	"context"
	"errors"

	"github.com/TomBowyerResearchProject/common/logger"
	commonMongo "github.com/TomBowyerResearchProject/common/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	PageLimit = 20
)

func FindUser(ctx context.Context, username string) (*model.User, error) {
	user := model.User{}
	filter := bson.D{primitive.E{Key: "username", Value: username}}

	db := commonMongo.GetDatabase()
	usersCollection := db.Collection(UsersCollection)
	err := usersCollection.FindOne(ctx, filter).Decode(&user)

	// Create the user.
	if errors.Is(err, mongo.ErrNoDocuments) {
		userdef, err := CreateUser(ctx, username)

		return userdef, err
	}

	return &user, err
}

func FindToken(ctx context.Context, token string) (*model.Token, error) {
	tokenObj := model.Token{}
	filter := bson.D{primitive.E{Key: "token", Value: token}}

	db := commonMongo.GetDatabase()
	tokenCollection := db.Collection(TokensCollection)
	err := tokenCollection.FindOne(ctx, filter).Decode(&tokenObj)

	// Create the user.
	if errors.Is(err, mongo.ErrNoDocuments) {
		return &tokenObj, messages.ErrNoToken
	}

	return &tokenObj, err
}

func GetAllUsers(ctx context.Context) *[]model.Connection {
	var userList []model.Connection

	db := commonMongo.GetDatabase()
	userCollection := db.Collection(UsersCollection)

	cursor, err := userCollection.Find(ctx, bson.D{})
	if errors.Is(err, mongo.ErrNoDocuments) {
		return &userList
	}

	for cursor.Next(ctx) {
		// Create a value into which the single document can be decoded.
		var connection model.Connection

		err := cursor.Decode(&connection)
		if err != nil {
			logger.Error(err)

			continue
		}

		userList = append(userList, connection)
	}

	return &userList
}

func GetMessagesForUsers(ctx context.Context, from, to string, skip int64) *[]model.ChatMessage {
	var chatList []model.ChatMessage

	fromQueryfrom := bson.M{"username_from": from}
	fromQueryTo := bson.M{"username_to": to}

	toQueryfrom := bson.M{"username_from": to}
	toQueryTo := bson.M{"username_to": from}

	query := bson.M{
		"$and": []bson.M{
			fromQueryfrom,
			fromQueryTo,
		},
	}

	secondQuery := bson.M{
		"$and": []bson.M{
			toQueryfrom,
			toQueryTo,
		},
	}

	full := bson.M{
		"$or": []bson.M{
			query,
			secondQuery,
		},
	}

	db := commonMongo.GetDatabase()
	messageCollection := db.Collection(MessageCollection)

	cursor, err := messageCollection.Find(ctx, full)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return &chatList
	}

	for cursor.Next(ctx) {
		// Create a value into which the single document can be decoded.
		var chatmessage model.ChatMessage

		err := cursor.Decode(&chatmessage)
		if err != nil {
			logger.Error(err)

			continue
		}

		chatList = append(chatList, chatmessage)
	}

	return &chatList
}
