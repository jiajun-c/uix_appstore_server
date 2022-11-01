package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
	"uix_appstore_server/handler"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	log.Println("Starting the application")

	router := mux.NewRouter()
	ctx, err := context.WithTimeout(context.Background(), 10*time.Second)
	if err != nil {
		fmt.Println(err)
	}
	var erro error
	handler.Client, erro = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if erro != nil {
		fmt.Println(erro)
	} else {
		fmt.Println("connect the server!")
	}
	router.HandleFunc("/api/user/login", handler.UserLogin).Methods("POST")
	router.HandleFunc("/api/user/signup", handler.UserSignup).Methods("POST")
	router.HandleFunc("/api/user/upload", handler.Upload)
	log.Fatal(http.ListenAndServe(":8080", router))

}
