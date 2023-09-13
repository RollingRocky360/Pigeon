package main

import (
	"Pigeon-Server/controller"
	"Pigeon-Server/model/database"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	database.InitDatabase()
	router := mux.NewRouter()
	wsServer := controller.NewWebsocketServer()

	router.HandleFunc("/chat", controller.CreateChat).Methods("POST")
	router.HandleFunc("/chat", controller.GetAllChatsOfUser).Methods("GET")
	router.HandleFunc("/chat/{chat_id}", controller.GetChatDetails).Methods("GET")

	router.HandleFunc("/user/{user_detail}", controller.GetUser).Methods("GET")

	router.HandleFunc("/messages/{chat_id}", controller.GetAllMessagesInChat).Methods("GET")
	router.HandleFunc("/ws/{user_id}", wsServer.HandleHandshake)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedHeaders:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "DELETE", "PUT", "OPTIONS"},
	})

	log.Println("Listening on :3000")
	log.Fatal(http.ListenAndServe(":3000", c.Handler(router)))
}
