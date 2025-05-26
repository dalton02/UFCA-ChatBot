package main

import (
	"licor_model/core/server"
	"licor_model/core/server/shared"
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Minha API
// @version 0.0.1
// @description Documentação do melhor chatbot do mundo2
// @host localhost:4000
// @BasePath /
func main() {

	go docs()

	database, _ := server.InitConnection()
	defer database.Close()
	shared.SetDB(database)
	server.MainServer()

}

func docs() {
	mux := http.NewServeMux()
	mux.Handle("/docs/", httpSwagger.WrapHandler)
	http.ListenAndServe(":4001", mux)
}
