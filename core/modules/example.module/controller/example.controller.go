package example_controller

import (
	example_service "licor_model/core/modules/example.module/service"
	"net/http"

	"github.com/dalton02/licor/httpkit"
)

func ExampleRoute(response http.ResponseWriter, request *http.Request) httpkit.HttpMessage {

	example_service.GetService()
	return httpkit.AppSucess("Welcome to licor :D", make(map[string]string))
}
