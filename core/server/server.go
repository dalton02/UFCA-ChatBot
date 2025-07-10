package server

import (
	chat_controller "licor_model/core/modules/chat/controller"
	chat_dto "licor_model/core/modules/chat/dto"
	scrapper_service "licor_model/core/modules/scrapper/service"
	_ "licor_model/docs"

	"github.com/dalton02/licor/licor"
)

func MainServer() {
	scrapper_service.Init()
	licor.Public[chat_dto.NovoChatDto, any]("/chat").Post(chat_controller.NovoChat)
	licor.Public[any, any]("/chat/{idChat}").Get(chat_controller.BuscarChat)
	licor.Public[chat_dto.NovaMensagemDto, any]("/chat/{idChat}/mensagem").Post(chat_controller.EnviarMensagem)
	licor.Init("3000")

}
