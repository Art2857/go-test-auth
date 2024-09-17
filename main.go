package main

import (
	"auth-service/pkg/config"
	"auth-service/pkg/database"
	"auth-service/pkg/token"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	config.Init()

	database.InitDB(config.Env.POSTGRES_CONNECTION)
	defer database.CloseDB()

	router := mux.NewRouter()
	token.InitRouter(router.PathPrefix("/token").Subrouter())

	log.Println("Server is running on port 8000...")
	log.Fatal(http.ListenAndServe(":8000", router))
}
