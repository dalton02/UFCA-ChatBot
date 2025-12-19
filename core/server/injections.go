package server

import (
	auth_controller "licor_model/core/modules/auth/controller"
	auth_middleware "licor_model/core/modules/auth/middleware"
	auth_service "licor_model/core/modules/auth/service"
	chat_controller "licor_model/core/modules/chat/controller"
	chat_middleware "licor_model/core/modules/chat/middleware"
	chat_service "licor_model/core/modules/chat/service"
	document_service "licor_model/core/modules/document/service"
	jsonbin_controller "licor_model/core/modules/jsonbin/controller"
	jsonbin_service "licor_model/core/modules/jsonbin/service"
	redis_service "licor_model/core/modules/redis"
)

func InitInjections() {

	//Services

	redisService, err := redis_service.NewRedisService()
	if err != nil {
		panic(err)
	}

	authService := auth_service.NewAuthService()
	docService := document_service.NewDocumentService()

	chatService := chat_service.NewChatService(docService)
	jsonService := jsonbin_service.NewJsonBinService(docService)

	//Middlewares
	authMid := auth_middleware.NewAuthMiddleware(authService, redisService)
	chatMid := chat_middleware.NewChatService(redisService)

	//Controllers
	chatControl := chat_controller.NewController(chatService, chatMid)
	authControl := auth_controller.NewAuthController(authService)
	jsonControl := jsonbin_controller.New(jsonService, authMid)

	//Agrupamento de rotas
	publicGroup := routes.Engine.Group("/", authMid.IpValidation)
	protectedGroup := routes.Engine.Group("/", authMid.JwtGuard)

	//Declaração de middlewares principais para agrupamento
	routes.Groups.PublicGroup = publicGroup
	routes.Groups.JwtGroup = protectedGroup

	chatControl.Routes(routes.Groups.JwtGroup)
	jsonControl.Routes(routes.Groups.JwtGroup)

	authControl.Routes(routes.Groups.PublicGroup)
}
