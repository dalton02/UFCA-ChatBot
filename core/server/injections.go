package server

import (
	chat_controller "licor_model/core/modules/chat/controller"
	chat_service "licor_model/core/modules/chat/service"
	document_service "licor_model/core/modules/document/service"
)

func InitInjections() {

	//Services

	docService := document_service.NewDocumentService()
	chatService := chat_service.NewChatService(docService)

	//Controllers
	chatControl := chat_controller.NewController(chatService)

	//Grupamento de rotas
	publicGroup := routes.Engine.Group("/")

	//Declaração de middlewares principais para agrupamento
	routes.Groups.PublicGroup = publicGroup

	//Inicialização de Rotas
	chatControl.Routes(routes.Groups.PublicGroup) //Aqui, as rotas de chats passam pelo middleware e grupo do JwtGroup

}
