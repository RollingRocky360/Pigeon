package controller

import (
	"Pigeon-Server/controller/utils"
	"Pigeon-Server/model"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateChat(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")

	var users []primitive.ObjectID
	err := json.NewDecoder(r.Body).Decode(&users)
	if err != nil {
		utils.WriteInternalServerError(&w, err)
		return
	}

	chat, err := model.CreateChat(users)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			utils.WriteInternalServerError(&w, err)
		}
		return
	}

	json.NewEncoder(w).Encode(chat)
}

func GetAllChatsOfUser(w http.ResponseWriter, r *http.Request) {
	userIDHex := strings.Split(r.Header.Get("Authorization"), " ")[1]

	chats, err := model.GetAllChatsOfUser(userIDHex)
	if err != nil {
		utils.WriteInternalServerError(&w, err)
		log.Println(err.Error())
		return
	}

	json.NewEncoder(w).Encode(chats)
}

func GetChatDetails(w http.ResponseWriter, r *http.Request) {
	chatId := mux.Vars(r)["chat_id"]

	chat, err := model.GetChatDetails(chatId)
	if err != nil {
		utils.WriteInternalServerError(&w, err)
		log.Println(err.Error())
		return
	}

	json.NewEncoder(w).Encode(chat)
}
