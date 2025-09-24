package chat_controller

import (
	chat_dto "licor_model/core/modules/chat/dto"
	chat_service "licor_model/core/modules/chat/service"
	util_dto "licor_model/core/util/dto"
	"licor_model/core/util/interceptor"
	"time"

	"github.com/gin-gonic/gin"
)

type ChatController struct {
	chatService *chat_service.ChatService
}

func NewController(chatService *chat_service.ChatService) *ChatController {
	return &ChatController{
		chatService: chatService,
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
	id, err := c.chatService.CreateChat(createChat)

	if err != nil {
		interceptor.AppBadRequest(ctx, err.Error())
		return
	}

	interceptor.AppSuccess(ctx, "Chat criado com sucesso", chat_dto.ChatDto{
		ID: id,
		TimeStampDefaultDB: util_dto.TimeStampDefaultDB{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		UserID: createChat.UserID,
		Title:  createChat.Title,
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

	response, err := c.chatService.SaveMessage(message, chatID)

	if err != nil {
		interceptor.AppBadRequest(ctx, err.Error())
		return
	}
	interceptor.AppSuccess(ctx, "Resposta da IA", response)
}

// GetChat godoc
// @Summary Buscar chat por ID
// @Description Retorna os detalhes de um chat específico
// @Tags Chats
// @Produce json
// @Param chatID path string true "ID do chat"
// @Success 200 {object} util_dto.AppResponse{data=chat_dto.ChatDto} "Chat encontrado"
// @Failure 404 {object} util_dto.AppResponse "Nenhum chat encontrado"
// @Router /chats/{chatID}/ [get]
func (c *ChatController) GetChat(ctx *gin.Context) {
	chatID, _ := ctx.Params.Get("chatID")
	result, err := c.chatService.GetChatByID(chatID)
	if err != nil {
		interceptor.AppNotFound(ctx, err.Error())
		return
	}
	interceptor.AppSuccess(ctx, "Chat encontrado", result)
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

	result, err := c.chatService.ListChat(filters)
	if err != nil {
		interceptor.AppNotFound(ctx, err.Error())
		return
	}
	interceptor.AppSuccess(ctx, "Chats encontrados", result)
}
func (c *ChatController) Routes(g *gin.RouterGroup) {

	g.POST("/chats/new-chat", c.CreateChat)
	g.GET("/chats/list-chats", c.ListChats)

	insideChat := g.Group("/chats/:chatID")
	{
		insideChat.GET("/", c.GetChat)
		insideChat.POST("/new-message", c.SendMessage)
	}
}
