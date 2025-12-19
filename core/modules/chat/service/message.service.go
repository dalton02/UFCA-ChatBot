package chat_service

import (
	"database/sql"
	"fmt"
	chat_dto "licor_model/core/modules/chat/dto"
	chat_repository "licor_model/core/modules/chat/repository"
	document_dto "licor_model/core/modules/document/dto"
	ollama_dto "licor_model/core/modules/ollama/dto"
	ai_service "licor_model/core/modules/ollama/service"
	"licor_model/core/util/executor"
	"licor_model/core/util/transaction"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func (s *ChatService) ListMessages(filters chat_dto.QueryListMessageDto, chatID string) (response chat_dto.ListMessageDto, err error) {
	return s.repo.ListMessages(filters, chatID)
}

func (s *ChatService) SaveMessage(ctx *gin.Context, message chat_dto.CreateMensagemDto, chatID string, origins []document_dto.OriginEnum) (response chat_dto.ResponseNewMessage, err error) {

	documents, err := s.docService.GetDocumentsBySimiliarity(message.Content, origins)
	if err != nil {
		return response, err
	}

	var links []string = []string{}
	for _, doc := range documents {
		links = append(links, doc.Link)
	}

	transaction.RunInTx(func(tx *sql.Tx) error {
		repoTransaction := chat_repository.NewChatRepository(executor.NewDBExecutor(tx))

		//Salvando mensagem do usuário
		idHuman, err := repoTransaction.CreateMessage(message.Content, chatID, false, []string{})
		if err != nil {
			return err
		}

		var stringBuilder strings.Builder
		stringBuilder.WriteString(`Você é um assistente da universidade. Responda com base nos dados abaixo e tente manter um foco maior apenas em duvidas da universidade no geral, tente sempre que possivel entrar em detalhes sobre a duvida que foi lhe mandada`)
		for _, doc := range documents {
			stringBuilder.WriteString("\n---\n")
			stringBuilder.WriteString(doc.Content)
			stringBuilder.WriteString("\n---")
		}
		stringBuilder.WriteString(`
		PERGUNTA: ` + message.Content + `
		RESPOSTA:`)
		fmt.Println(stringBuilder.String())
		responseIA, err := ai_service.SendRequestOpenAIStream(ctx, links, ollama_dto.RequestChatAI{
			Model: os.Getenv("OLLAMA_MODEL"),
			Messages: []ollama_dto.MessageChatAI{
				{
					Content: stringBuilder.String(),
					Role:    "user",
				},
			},
		})
		if err != nil {
			return err
		}

		idIA, err := s.repo.CreateMessage(responseIA.Message.Content, chatID, true, links)
		if err != nil {
			return err
		}

		response.IDIa = idIA
		response.IDHuman = idHuman
		response.IAContent = responseIA.Message.Content
		response.Links = links
		return nil

	})

	return response, err
}
