package main

import (
	"fmt"
	"log"
	"net/http" // missing in your code
	"os"

	handler "project/Handler"
	"project/db"
	"project/middleware"
	"project/repository"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	// Load env
	godotenv.Load()

	mongoURI := os.Getenv("MONGO_URI")
	dbName := os.Getenv("MONGO_DB")

	// Connect to Mongo
	client, database, err := db.ConnectMongo(mongoURI, dbName)
	if err != nil {
		log.Fatal("Mongo connection failed:", err)
	}

	// Init collection inside repository
	repository.InitCollections(client, database)

	// Router
	r := mux.NewRouter()
	r.HandleFunc("/account", handler.CreateAccount).Methods("POST")
	r.HandleFunc("/account/update/{id}", handler.UpdateAccountHandler).Methods("PUT")
	r.HandleFunc("/account/get/{id}", handler.FetchAccountHandler).Methods("GET")
	r.HandleFunc("/account/delete/{id}", handler.DeleteAccountHandler).Methods("DELETE")
	r.HandleFunc("/entry/update/{id}", handler.UpdateEntryHandler).Methods("PUT")
	r.HandleFunc("/entry", handler.CreateEntryHandler).Methods("POST")
	r.HandleFunc("/transfer/update/{id}", handler.UpdateTransferHandler).Methods("PUT")
	r.HandleFunc("/transfer", handler.InsertTransferHandler).Methods("POST")
	r.HandleFunc("/user", handler.InsertUserHandler).Methods("POST")
	r.HandleFunc("/user/login", handler.LoginHandler).Methods("POST")
	r.HandleFunc("/users", middleware.Protect(middleware.AdminOnly(handler.FetchAllUserHandler))).Methods("GET")
	r.HandleFunc("/transfer/account/{id}", middleware.Protect(handler.FetchMyTransfersHandler)).Methods("GET")
	r.HandleFunc("/account/balance/{id}", middleware.Protect(handler.GetBalanceHandler)).Methods("GET")
	fmt.Println("🚀 Server started at :4000")
	log.Fatal(http.ListenAndServe(":4000", middleware.EnableCORS(r)))

}
