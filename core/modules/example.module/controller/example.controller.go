package example_controller

import (
	example_service "licor_model/core/modules/example.module/service"
	"net/http"

	"github.com/dalton02/licor/httpkit"
)

// @Summary Rota de exemplo licor :D
// @Description Obtém um hello world mt bacana
// @Tags Tag de Exemplos
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /hello-world [get]
func ExampleRoute(response http.ResponseWriter, request *http.Request) httpkit.HttpMessage {

	example_service.GetService()
	return httpkit.AppSucess("Welcome to licor :D", make(map[string]string))
}
