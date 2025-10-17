package server

import (
	auth_controller "licor_model/core/modules/auth/controller"
	auth_middleware "licor_model/core/modules/auth/middleware"
	auth_service "licor_model/core/modules/auth/service"
	chat_controller "licor_model/core/modules/chat/controller"
	chat_service "licor_model/core/modules/chat/service"
	document_service "licor_model/core/modules/document/service"
	scrapper_service "licor_model/core/modules/scrapper/service"
)

func InitInjections() {

	//Services
	docService := document_service.NewDocumentService()
	chatService := chat_service.NewChatService(docService)
	scrapperService := scrapper_service.NewScrapperService(docService)
	authService := auth_service.NewAuthService()

	//Middlewares
	authMid := auth_middleware.NewAuthMiddleware(authService)

	//Controllers
	chatControl := chat_controller.NewController(chatService)
	authControl := auth_controller.NewAuthController(authService)

	//Agrupamento de rotas
	publicGroup := routes.Engine.Group("/", authMid.IpValidation)
	protectedGroup := routes.Engine.Group("/", authMid.JwtGuard)

	//Declaração de middlewares principais para agrupamento
	routes.Groups.PublicGroup = publicGroup
	routes.Groups.JwtGroup = protectedGroup

	chatControl.Routes(routes.Groups.PublicGroup)
	authControl.Routes(routes.Groups.JwtGroup)
	scrapperService.Init()
}
