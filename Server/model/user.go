package model

import (
	"Pigeon-Server/model/database"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID       primitive.ObjectID   `json:"id" bson:"_id"`
	Email    string               `json:"email"`
	Username string               `json:"username"`
	Password string               `json:"password"`
	Chats    []primitive.ObjectID `json:"chats"`
}

func SearchUser(userDetail string) (User, error) {

	var user User
	err := database.GetCollection("user").FindOne(context.TODO(), bson.D{{
		Key: "$or",
		Value: []bson.D{
			{{Key: "username", Value: userDetail}},
			{{Key: "email", Value: userDetail}},
		},
	}}).Decode(&user)

	return user, err
}
