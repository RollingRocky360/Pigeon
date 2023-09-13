package model

import (
	"Pigeon-Server/model/database"
	"context"
	"log"
	"sync"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Message struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	ChatID    primitive.ObjectID `json:"chat_id" bson:"chat_id"`
	UserID    primitive.ObjectID `json:"user_id" bson:"user_id"`
	Content   string             `json:"content"`
	Timestamp int64              `json:"timestamp"`
}

type messageBuffer struct {
	Buffer  []interface{}
	MaxSize int
	Mutex   sync.Mutex
}

func NewMessageBuffer() *messageBuffer {
	return &messageBuffer{MaxSize: 5}
}

func (mb *messageBuffer) Insert(m Message) {
	mb.Mutex.Lock()
	mb.Buffer = append(mb.Buffer, m)
	mb.Mutex.Unlock()

	if len(mb.Buffer) != mb.MaxSize {
		return
	}

	mb.flush()
}

func (mb *messageBuffer) flush() {
	if len(mb.Buffer) == 0 {
		return
	}

	Messages := database.GetCollection("message")

	log.Println("flushing")
	mb.Mutex.Lock()
	_, err := Messages.InsertMany(context.TODO(), mb.Buffer)
	mb.Mutex.Unlock()

	if err != nil {
		panic(err)
	}

	mb.Buffer = nil
}

var BUF = NewMessageBuffer()

func GetMessagesInChat(chat_id primitive.ObjectID) ([]Message, error) {
	BUF.flush()

	Messages := database.GetCollection("message")

	cursor, err := Messages.Find(context.TODO(), bson.D{{Key: "chat_id", Value: chat_id}})
	if err != nil {
		return []Message{}, err
	}

	messages := []Message{}

	if err = cursor.All(context.TODO(), &messages); err != nil {
		return messages, nil
	}

	return messages, nil
}

func DeleteMessage(m Message) error {
	BUF.flush()
	_, err := database.GetCollection("message").DeleteOne(context.TODO(), bson.D{{Key: "_id", Value: m.ID}})
	return err
}

func InsertMessage(m Message) {
	BUF.Insert(m)
}
