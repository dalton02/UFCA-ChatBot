package server

import (
	scrapper_service "licor_model/core/modules/scrapper/service"
	_ "licor_model/docs"

	"github.com/dalton02/licor/licor"
)

func MainServer() {

	scrapper_service.Init()
	licor.Init("4000")

}
