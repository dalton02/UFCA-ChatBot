package chat_service

import (
	"database/sql"
	chat_dto "licor_model/core/modules/chat/dto"
	chat_repository "licor_model/core/modules/chat/repository"
	ollama_dto "licor_model/core/modules/ollama/dto"
	ollama_service "licor_model/core/modules/ollama/service"
	"licor_model/core/util/executor"
	"licor_model/core/util/transaction"
	"os"
	"strings"
)

func (s *ChatService) ListMessages(filters chat_dto.QueryListMessageDto, chatID string) (response chat_dto.ListMessageDto, err error) {
	return s.repo.ListMessages(filters, chatID)
}

func (s *ChatService) SaveMessage(message chat_dto.CreateMensagemDto, chatID string) (response chat_dto.ResponseNewMessage, err error) {

	documents, err := s.docService.GetDocumentsBySimiliarity(message.Content)
	if err != nil {
		return response, err
	}

	transaction.RunInTx(func(tx *sql.Tx) error {
		repoTransaction := chat_repository.NewChatRepository(executor.NewDBExecutor(tx))

		//Salvando mensagem do usuário
		idHuman, err := repoTransaction.CreateMessage(message.Content, chatID, false)
		if err != nil {
			return err
		}

		var stringBuilder strings.Builder
		stringBuilder.WriteString(`Você é um assistente da universidade. Responda APENAS com base nos dados abaixo. Nunca mencione sites. Mantenha o foco em alunos e universidade.
		DADOS:`)
		for _, doc := range documents {
			stringBuilder.WriteString("\n---\n")
			stringBuilder.WriteString(doc.Content)
			stringBuilder.WriteString("\n---")
		}
		stringBuilder.WriteString(`
		PERGUNTA: ` + message.Content + `
		RESPOSTA:`)
		responseIA, err := ollama_service.SendRequest(ollama_dto.RequestChatAI{
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

		idIA, err := s.repo.CreateMessage(responseIA.Message.Content, chatID, true)
		if err != nil {
			return err
		}

		response.IDIa = idIA
		response.IDHuman = idHuman
		response.IAContent = responseIA.Message.Content
		return nil

	})

	return response, err
}
