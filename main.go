package main

import (
	"fmt"
	"licor_model/core/server"
	"licor_model/core/server/shared"
)

// @title Golang CHATBOT
// @version @0.0.1
// @description Documentação do melhor chatbot universitario
// @BasePath /
// @SecurityDefinitions.apikey Bearer Auth
// @in header
// @name Authorization
func main() {

	database, err := server.InitConnection()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer database.Close()
	shared.SetDB(database)
	server.MainServer()

}
