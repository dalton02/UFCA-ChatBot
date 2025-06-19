package server

import (
	chat_controller "licor_model/core/modules/chat/controller"
	_ "licor_model/docs"

	"github.com/dalton02/licor/licor"
)

func MainServer() {

	licor.Public[any, any]("/chat/mensagem").Post(chat_controller.EnviarMensagem)
	licor.Init("3000")

}
