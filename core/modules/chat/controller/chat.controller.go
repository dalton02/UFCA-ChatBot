package chat_controller

import (
	"net/http"

	"github.com/dalton02/licor/httpkit"
)

func EnviarMensagem(response http.ResponseWriter, request *http.Request) httpkit.HttpMessage {

	return httpkit.AppSuccess("Tudo certo aqui!", make(map[string]interface{}))
}
