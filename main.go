package main

import (
	"licor_model/core/server"
	"licor_model/core/server/shared"
)

func main() {

	database, _ := server.InitConnection()
	defer database.Close()
	shared.SetDB(database)
	server.MainServer()

}
