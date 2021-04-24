package db

import (
	"chatter/internal/messages"
	"chatter/model"
	"context"

	"github.com/TomBowyerResearchProject/common/logger"
	commonMongo "github.com/TomBowyerResearchProject/common/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	PageLimit = 20
)

func FindUser(username string) (*model.User, error) {
	user := model.User{}
	filter := bson.D{primitive.E{Key: "username", Value: username}}

	db := commonMongo.GetDatabase()
	usersCollection := db.Collection(UsersCollection)
	err := usersCollection.FindOne(context.TODO(), filter).Decode(&user)

	// Create the user
	if err == mongo.ErrNoDocuments {
		userdef, err := CreateUser(username)
		return userdef, err
	}
	return &user, err
}

func FindToken(token string) (*model.Token, error) {
	tokenObj := model.Token{}
	filter := bson.D{primitive.E{Key: "token", Value: token}}

	db := commonMongo.GetDatabase()
	tokenCollection := db.Collection(TokensCollection)
	err := tokenCollection.FindOne(context.TODO(), filter).Decode(&tokenObj)

	// Create the user
	if err == mongo.ErrNoDocuments {
		return &tokenObj, messages.ErrNoToken
	}
	return &tokenObj, err
}

func GetAllUsers() *[]model.Connection {
	var userList []model.Connection

	db := commonMongo.GetDatabase()
	userCollection := db.Collection(UsersCollection)
	cursor, err := userCollection.Find(context.TODO(), bson.D{})
	if err == mongo.ErrNoDocuments {
		return &userList
	}

	for cursor.Next(context.TODO()) {
		//Create a value into which the single document can be decoded
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

func GetMessagesForUsers(from, to string, skip int64) *[]model.ChatMessage {
	var chatList []model.ChatMessage

	query := bson.M{
		"$and": []bson.M{
			bson.M{"username_from": from},
			bson.M{"username_to": to},
		},
	}

	secondQuery := bson.M{
		"$and": []bson.M{
			bson.M{"username_from": to},
			bson.M{"username_to": from},
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
	cursor, err := messageCollection.Find(context.TODO(), full)
	if err == mongo.ErrNoDocuments {
		return &chatList
	}

	for cursor.Next(context.TODO()) {
		//Create a value into which the single document can be decoded
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
