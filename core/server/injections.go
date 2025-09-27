package server

import (
	"database/sql"
	auth_controller "licor_model/core/modules/auth/controller"
	auth_repository "licor_model/core/modules/auth/repository"
	auth_service "licor_model/core/modules/auth/service"
	chat_controller "licor_model/core/modules/chat/controller"
	chat_service "licor_model/core/modules/chat/service"
	document_service "licor_model/core/modules/document/service"
	scrapper_service "licor_model/core/modules/scrapper/service"
	"licor_model/core/util/executor"
)

func InitInjections() {

	//não sei ao certo se a variável que tem que ser passada em exec é essa tx genérica mesmo
	var tx *sql.Tx
	exec := executor.NewDBExecutor(tx)
	authRepo := auth_repository.NewAuthRepository(exec)

	//Services
	docService := document_service.NewDocumentService()
	chatService := chat_service.NewChatService(docService)
	scrapperService := scrapper_service.NewScrapperService(docService)
	authService := auth_service.NewAuthService(authRepo)

	//Controllers
	chatControl := chat_controller.NewController(chatService)
	authControl := auth_controller.NewAuthController(authService)

	//Agrupamento de rotas
	publicGroup := routes.Engine.Group("/")

	//Declaração de middlewares principais para agrupamento
	routes.Groups.PublicGroup = publicGroup

	//Inicialização de Rotas

	chatControl.Routes(routes.Groups.PublicGroup) //Aqui, as rotas de chats passam pelo middleware e grupo do JwtGroup
	authControl.Routes(routes.Groups.PublicGroup) //Aqui, as rotas de auth não passam pelo middleware do JwtGroup

	scrapperService.Init()
}
