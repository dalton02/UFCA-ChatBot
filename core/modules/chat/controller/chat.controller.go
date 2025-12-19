package chat_controller

import (
	"fmt"
	chat_dto "licor_model/core/modules/chat/dto"
	chat_middleware "licor_model/core/modules/chat/middleware"
	chat_service "licor_model/core/modules/chat/service"
	document_dto "licor_model/core/modules/document/dto"
	util_dto "licor_model/core/util/dto"
	"licor_model/core/util/interceptor"
	"time"

	"github.com/gin-gonic/gin"
)

type ChatController struct {
	chatService    *chat_service.ChatService
	chatMiddleware *chat_middleware.ChatMiddleware
}

func NewController(chatService *chat_service.ChatService, chatMiddleware *chat_middleware.ChatMiddleware) *ChatController {
	return &ChatController{
		chatService:    chatService,
		chatMiddleware: chatMiddleware,
	}
}

// CreateChat godoc
// @Summary Criar novo chat
// @Description Cria um novo chat com título e usuário associado
// @Tags Chats
// @Accept json
// @Produce json
// @Param request body chat_dto.CreateChatDto true "Dados para criar chat"
// @Success 200 {object} util_dto.AppResponse{data=chat_dto.ChatDto} "Chat criado com sucesso"
// @Failure 400 {object} util_dto.AppResponse "Erro de validação ou requisição"
// @Router /chats/new-chat [post]
func (c *ChatController) CreateChat(ctx *gin.Context) {
	var createChat chat_dto.CreateChatDto

	if err := interceptor.ValidateAndExtract(ctx, &createChat); err != nil {
		interceptor.AppBadRequest(ctx, err.Error())
		return
	}

	userID := ctx.MustGet("userID").(string)

	id, err := c.chatService.CreateChat(createChat, userID)

	if err != nil {
		interceptor.AppBadRequest(ctx, err.Error())
		return
	}

	interceptor.AppSuccess(ctx, "Chat criado com sucesso", chat_dto.ChatDto{
		ID:     id,
		UserID: userID,
		TimeStampDefaultDB: util_dto.TimeStampDefaultDB{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Title: createChat.Title,
	})

}

// SendMessage godoc
// @Summary Enviar mensagem para um chat
// @Description Envia uma nova mensagem dentro de um chat existente
// @Tags Chats
// @Accept json
// @Produce json
// @Param chatID path string true "ID do chat"
// @Param request body chat_dto.CreateMensagemDto true "Dados da mensagem"
// @Success 200 {object} util_dto.AppResponse{data=string} "Resposta da IA"
// @Failure 400 {object} util_dto.AppResponse "Erro de validação ou requisição"
// @Router /chats/{chatID}/new-message [post]
func (c *ChatController) SendMessage(ctx *gin.Context) {
	var message chat_dto.CreateMensagemDto

	if err := interceptor.ValidateAndExtract(ctx, &message); err != nil {
		interceptor.AppBadRequest(ctx, err.Error())
		return
	}

	chatID, _ := ctx.Params.Get("chatID")

	response, err := c.chatService.SaveMessage(ctx, message, chatID, []document_dto.OriginEnum{document_dto.Wikiversity, document_dto.JsonStudents})

	if err != nil {
		interceptor.AppBadRequest(ctx, err.Error())
		return
	}
	interceptor.AppSuccess(ctx, "Resposta da IA", response)
}

// ListChats godoc
// @Summary Listar chats
// @Description Lista chats aplicando filtros opcionais
// @Tags Chats
// @Accept json
// @Produce json
// @Param query query chat_dto.ListChatDto true "Filtros"
// @Success 200 {object} util_dto.AppResponse{data=chat_dto.ListChatDto} "Lista de chats"
// @Failure 404 {object} util_dto.AppResponse "Nenhum chat encontrado"
// @Router /chats/list-chats [get]
func (c *ChatController) ListChats(ctx *gin.Context) {
	var filters chat_dto.QueryListChatDto

	if err := interceptor.ValidateAndExtractQuery(ctx, &filters); err != nil {
		interceptor.AppBadRequest(ctx, err.Error())
		return
	}

	userID := ctx.MustGet("userID").(string)

	result, err := c.chatService.ListChat(filters, userID)
	if err != nil {
		interceptor.AppError(ctx, err)
		return
	}
	interceptor.AppSuccess(ctx, "Chats encontrados", result)
}

func (c *ChatController) DeleteChat(ctx *gin.Context) {

	chatID := ctx.Param("chatID")

	userID := ctx.MustGet("userID").(string)

	err := c.chatService.DeleteChat(chatID, userID)
	if err != nil {
		interceptor.AppError(ctx, err)
		return
	}
	interceptor.AppSuccess(ctx, "Chats Deletado", nil)
}

// ListChats godoc
// @Summary Listar mensagens
// @Tags Chats
// @Accept json
// @Produce json
// @Param query query chat_dto.ListMessageDto true "Filtros"
// @Success 200 {object} util_dto.AppResponse{data=chat_dto.ListMessageDto} "Lista de chats"
// @Failure 404 {object} util_dto.AppResponse "Nenhum chat encontrado"
// @Router /chats/{chatID}/list-messages [get]
func (c *ChatController) ListMessages(ctx *gin.Context) {
	var filters chat_dto.QueryListMessageDto

	if err := interceptor.ValidateAndExtractQuery(ctx, &filters); err != nil {
		interceptor.AppBadRequest(ctx, err.Error())
		return
	}

	chatID, _ := ctx.Params.Get("chatID")

	result, err := c.chatService.ListMessages(filters, chatID)
	fmt.Println(result, err)
	if err != nil {
		interceptor.AppError(ctx, err)
		return
	}
	interceptor.AppSuccess(ctx, "Chats encontrados", result)
}

func (c *ChatController) Routes(g *gin.RouterGroup) {

	g.POST("/chats/new-chat", c.CreateChat)
	g.GET("/chats/list-chats", c.ListChats)

	insideChat := g.Group("/chats/:chatID")
	{
		insideChat.DELETE("/", c.DeleteChat)
		insideChat.GET("/list-messages", c.ListMessages)
		insideChat.POST("/new-message", c.chatMiddleware.ValidateAIDailyUse, c.SendMessage)
	}
}
