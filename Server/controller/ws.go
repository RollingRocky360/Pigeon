package controller

import (
	"Pigeon-Server/model"
	"Pigeon-Server/model/database"
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type websocketServer struct {
	liveChats map[string]map[*websocket.Conn]bool
	mu        sync.Mutex
}

type websocketMessage struct {
	Type    string        `json:"type"`
	Message model.Message `json:"message"`
}

func NewWebsocketServer() *websocketServer {
	return &websocketServer{liveChats: make(map[string]map[*websocket.Conn]bool)}
}

func (s *websocketServer) addUserChats(conn *websocket.Conn, user_id string) error {
	var u model.User
	Users := database.GetCollection("user")

	uid, _ := primitive.ObjectIDFromHex(user_id)
	err := Users.FindOne(context.TODO(), bson.D{{Key: "_id", Value: uid}}).Decode(&u)
	if err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	for _, chat := range u.Chats {
		if s.liveChats[chat.Hex()] == nil {
			s.liveChats[chat.Hex()] = make(map[*websocket.Conn]bool)
		}
		s.liveChats[chat.Hex()][conn] = true
	}

	return nil
}

func (s *websocketServer) RemoveUserChats(conn *websocket.Conn, user_id string) error {
	var u model.User
	Users := database.GetCollection("user")

	uid, _ := primitive.ObjectIDFromHex(user_id)
	err := Users.FindOne(context.TODO(), bson.D{{Key: "_id", Value: uid}}).Decode(&u)
	if err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	for _, chat := range u.Chats {
		hexValue := chat.Hex()
		delete(s.liveChats[hexValue], conn)
		if s.liveChats[hexValue] == nil {
			delete(s.liveChats, hexValue)
		}
	}

	return nil
}

func (s *websocketServer) HandleHandshake(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		fmt.Println("websocket: connection failed -", err)
		return
	}

	user_id := mux.Vars(r)["user_id"]
	err = s.addUserChats(conn, user_id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			w.Write([]byte(""))
			return
		}
		log.Println(err.Error())
	}

	var m *model.Message
	wsm := websocketMessage{}
	for {
		if err := conn.ReadJSON(&wsm); err != nil {
			log.Println("websocket: User disconnected")
			s.RemoveUserChats(conn, user_id)

			conn.Close()
			break
		}

		m = &wsm.Message

		switch wsm.Type {
		case "new":
			m.ID = primitive.NewObjectID()
			m.Timestamp = time.Now().UnixNano()
			log.Println("Message:", m.Content)
			model.InsertMessage(*m)
		case "delete":
			err := model.DeleteMessage(*m)
			if err != nil {
				conn.WriteMessage(websocket.TextMessage, []byte("Error"))
				continue
			}
			log.Println("Deleted:", m.Content)
		}

		for c := range s.liveChats[m.ChatID.Hex()] {
			c.WriteJSON(wsm)
		}
	}
}
