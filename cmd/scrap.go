package cmd

import (
	document_service "licor_model/core/modules/document/service"
	scrapper_service "licor_model/core/modules/scrapper/service"
)

func main() {

	docService := document_service.NewDocumentService()

	scrapperService := scrapper_service.NewScrapperService(docService)
	scrapperService.Init()

}
