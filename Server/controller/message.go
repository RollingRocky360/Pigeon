package controller

import (
	"Pigeon-Server/controller/utils"
	"Pigeon-Server/model"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetAllMessagesInChat(w http.ResponseWriter, r *http.Request) {
	chatId := mux.Vars(r)["chat_id"]
	id, _ := primitive.ObjectIDFromHex(chatId)
	messages, err := model.GetMessagesInChat(id)

	if err != nil {
		utils.WriteInternalServerError(&w, err)
		log.Println(err.Error())
		return
	}

	json.NewEncoder(w).Encode(messages)
}
