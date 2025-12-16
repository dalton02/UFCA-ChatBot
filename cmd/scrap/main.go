package main

import (
	"fmt"
	document_service "licor_model/core/modules/document/service"
	scrapper_service "licor_model/core/modules/scrapper/service"
	"licor_model/core/server"
	"licor_model/core/server/shared"
)

func main() {
	database, err := server.InitConnection()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer database.Close()
	shared.SetDB(database)
	docService := document_service.NewDocumentService()

	scrapperService := scrapper_service.NewScrapperService(docService)
	scrapperService.Init()

}
