package server

import (
	auth_controller "licor_model/core/modules/auth/controller"
	auth_middleware "licor_model/core/modules/auth/middleware"
	auth_service "licor_model/core/modules/auth/service"
	chat_controller "licor_model/core/modules/chat/controller"
	chat_middleware "licor_model/core/modules/chat/middleware"
	chat_service "licor_model/core/modules/chat/service"
	document_service "licor_model/core/modules/document/service"
	redis_service "licor_model/core/modules/redis"
)

func InitInjections() {

	//Services

	redisService, err := redis_service.NewRedisService()
	if err != nil {
		panic(err)
	}

	docService := document_service.NewDocumentService()
	chatService := chat_service.NewChatService(docService)
	authService := auth_service.NewAuthService()

	//Middlewares
	authMid := auth_middleware.NewAuthMiddleware(authService)
	chatMid := chat_middleware.NewChatService(redisService)

	//Controllers
	chatControl := chat_controller.NewController(chatService, chatMid)
	authControl := auth_controller.NewAuthController(authService)

	//Agrupamento de rotas
	publicGroup := routes.Engine.Group("/", authMid.IpValidation)
	protectedGroup := routes.Engine.Group("/", authMid.JwtGuard)

	//Declaração de middlewares principais para agrupamento
	routes.Groups.PublicGroup = publicGroup
	routes.Groups.JwtGroup = protectedGroup

	chatControl.Routes(routes.Groups.JwtGroup)
	authControl.Routes(routes.Groups.PublicGroup)

}
