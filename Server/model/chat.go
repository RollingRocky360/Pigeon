package model

import (
	"Pigeon-Server/model/database"
	"context"
	"math/rand"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var pastels = []string{
	"#77DD77", "#836953", "#89cff0", "#99c5c4", "#9adedb",
	"#aa9499", "#aaf0d1", "#b2fba5", "#b39eb5", "#bdb0d0",
	"#bee7a5", "#befd73", "#c1c6fc", "#c6a4a4", "#cb99c9",
	"#ff6961", "#ff694f", "#ff9899", "#ffb7ce", "#ca9bf7",
}

type Chat struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	IsGroup   bool               `json:"is_group" bson:"is_group"`
	GroupName string             `json:"group_name" bson:"group_name"`
	Users     []User             `json:"users" bson:"users"`
	PFPColor  string             `json:"pfp_color" bson:"pfp_color"`
}

func CreateChat(users []primitive.ObjectID) (Chat, error) {
	Chats := database.GetCollection("chat")
	Users := database.GetCollection("user")

	u := []User{}

	userFindOp := bson.D{{Key: "_id", Value: bson.D{{Key: "$in", Value: users}}}}
	cursor, err := Users.Find(context.TODO(), userFindOp)
	if err != nil {
		return Chat{}, err
	}

	if err = cursor.All(context.TODO(), &u); err != nil {
		return Chat{}, err
	}

	chat := Chat{
		ID:       primitive.NewObjectID(),
		Users:    u,
		PFPColor: pastels[rand.Int()%len(pastels)],
	}

	for i := 0; i < len(chat.Users); i++ {
		chat.Users[i].Chats = append(chat.Users[i].Chats, chat.ID)
	}

	_, err = Chats.InsertOne(context.TODO(), chat)
	if err != nil {
		return Chat{}, err
	}

	if err != nil {
		return Chat{}, err
	}

	_, err = Users.UpdateMany(context.TODO(), userFindOp, bson.D{{Key: "$push", Value: bson.D{{Key: "chats", Value: chat.ID}}}})

	if err != nil {
		return chat, err
	}

	return chat, nil
}

func GetAllChatsOfUser(userIDHex string) ([]Chat, error) {

	userID, _ := primitive.ObjectIDFromHex(userIDHex)

	var user User
	err := database.GetCollection("user").FindOne(context.TODO(), bson.D{{Key: "_id", Value: userID}}).Decode(&user)
	if err != nil {
		return []Chat{}, err
	}

	var chats []Chat

	filter := bson.D{{Key: "_id", Value: bson.D{{Key: "$in", Value: user.Chats}}}}
	cursor, err := database.GetCollection("chat").Find(context.TODO(), filter)
	if err != nil {
		return chats, err
	}

	if err = cursor.All(context.TODO(), &chats); err != nil {
		return chats, err
	}
	return chats, nil
}

func GetChatDetails(chatHexID string) (Chat, error) {
	Chats := database.GetCollection("chat")

	chatID, err := primitive.ObjectIDFromHex(chatHexID)
	if err != nil {
		return Chat{}, err
	}

	var chat Chat
	err = Chats.FindOne(context.TODO(), bson.D{{Key: "_id", Value: chatID}}).Decode(&chat)
	if err != nil {
		return chat, err
	}

	return chat, nil
}
