package main

import (
	"fmt"
	"licor_model/core/server"
	"licor_model/core/server/shared"
	"net/http"

	"github.com/MarceloPetrucio/go-scalar-api-reference"
)

// @title Minha API
// @version 0.0.1
// @description Documentação do melhor chatbot do mundo2
// @host localhost:4000
// @BasePath /
func main() {

	database, _ := server.InitConnection()
	defer database.Close()
	shared.SetDB(database)

	http.HandleFunc("/docs", func(w http.ResponseWriter, r *http.Request) {
		htmlContent, err := scalar.ApiReferenceHTML(&scalar.Options{
			SpecURL: "./docs/swagger.json",
			CustomOptions: scalar.CustomOptions{
				PageTitle: "Chatbot swagger",
			},
			DarkMode: true,
			Layout:   scalar.LayoutModern,
			Theme:    scalar.ThemeMoon,
		})

		if err != nil {
			fmt.Printf("%v", err)
		}

		fmt.Fprintln(w, htmlContent)
	})

	server.MainServer()

}
