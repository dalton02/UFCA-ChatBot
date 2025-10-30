package chat_service

import (
	"database/sql"
	"fmt"
	chat_dto "licor_model/core/modules/chat/dto"
	chat_repository "licor_model/core/modules/chat/repository"
	ollama_dto "licor_model/core/modules/ollama/dto"
	ollama_service "licor_model/core/modules/ollama/service"
	"licor_model/core/util/executor"
	"licor_model/core/util/transaction"
	"os"
	"strings"
	"time"
)

func (s *ChatService) ListMessages(filters chat_dto.QueryListMessageDto, chatID string) (response chat_dto.ListMessageDto, err error) {
	return s.repo.ListMessages(filters, chatID)
}

func (s *ChatService) SaveMessage(message chat_dto.CreateMensagemDto, chatID string) (response chat_dto.MessageDto, err error) {

	documents, err := s.docService.GetDocumentsBySimiliarity(message.Content)
	if err != nil {
		return response, err
	}

	transaction.RunInTx(func(tx *sql.Tx) error {
		repoTransaction := chat_repository.NewChatRepository(executor.NewDBExecutor(tx))

		//Salvando mensagem do usuário
		_, err = repoTransaction.CreateMessage(message.Content, chatID, false)
		if err != nil {
			return err
		}

		var stringBuilder strings.Builder
		stringBuilder.WriteString("REGRAS DO PROMPT:")
		stringBuilder.WriteString("\n1 - Nunca envie links para websites em sua resposta ")
		stringBuilder.WriteString("\n2 - Sempre mantenha o foco do assunto a respeito da universidade e alunos")
		stringBuilder.WriteString("\n3 - Responda a pergunta baseando-se nos seguintes dados: ")
		stringBuilder.WriteString("\n----------INICIO DOS DADOS-----------")
		for _, doc := range documents {
			stringBuilder.WriteString("\n-----------------------------")
			stringBuilder.WriteString(doc.Content)
			stringBuilder.WriteString("\n-----------------------------")
		}
		stringBuilder.WriteString("\n------ FIM DOS DADOS-------")
		stringBuilder.WriteString("\nPergunta do usuário: ")
		stringBuilder.WriteString(message.Content)
		fmt.Println(stringBuilder.String())
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

		//Salvando mensagem da IA
		id, err := s.repo.CreateMessage(responseIA.Message.Content, chatID, true)

		response = chat_dto.MessageDto{
			ID:        id,
			ChatID:    chatID,
			Content:   responseIA.Message.Content,
			CreatedAt: time.Now(),
			Assistant: true,
		}

		return nil

	})

	return response, err
}
