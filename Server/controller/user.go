package controller

import (
	"Pigeon-Server/controller/utils"
	"Pigeon-Server/model"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetUser(w http.ResponseWriter, r *http.Request) {
	userDetail := mux.Vars(r)["user_detail"]

	var user model.User
	user, err := model.SearchUser(userDetail)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			w.Write([]byte("null"))
		} else {
			utils.WriteInternalServerError(&w, err)
		}
		return
	}

	json.NewEncoder(w).Encode(user)
}
