package chat_service

import (
	"fmt"
	chat_dto "licor_model/core/modules/chat/dto"
	chat_repository "licor_model/core/modules/chat/repository"
	document_service "licor_model/core/modules/document/service"
	ollama_dto "licor_model/core/modules/ollama/dto"
	ollama_service "licor_model/core/modules/ollama/service"
	"licor_model/core/util/executor"
	"os"
	"strings"
	"time"
)

type ChatService struct {
	repo       *chat_repository.ChatRepository
	docService *document_service.DocumentService
}

func NewChatService(docsService *document_service.DocumentService) *ChatService {
	return &ChatService{
		repo:       chat_repository.NewChatRepository(executor.NewDBExecutor(nil)),
		docService: docsService,
	}
}

func (s *ChatService) CreateChat(chat chat_dto.CreateChatDto, userID string) (id string, err error) {
	return s.repo.CreateChat(chat, userID)
}

func (s *ChatService) ListChat(filters chat_dto.QueryListChatDto) (chat_dto.ListChatDto, error) {
	return s.repo.ListChat(filters)
}

func (s *ChatService) GetChatByID(id string) (chat_dto.ChatDto, error) {
	return s.repo.GetChatByID(id)
}

func (s *ChatService) GetMessages(id int) ([]string, error) {
	return s.repo.GetMessages(id)
}

func (s *ChatService) SaveMessage(message chat_dto.CreateMensagemDto, chatID string) (response chat_dto.MessageDto, err error) {

	documents, err := s.docService.GetDocumentsBySimiliarity(message.Content)
	if err != nil {
		return response, err
	}
	//Salvando mensagem do usuário
	_, err = s.repo.CreateMessage(message.Content, chatID, false)

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
		return response, err
	}

	//Salvando mensagem da IA
	id, err := s.repo.CreateMessage(responseIA.Message.Content, chatID, true)

	return chat_dto.MessageDto{
		ID:        id,
		ChatID:    chatID,
		Content:   responseIA.Message.Content,
		CreatedAt: time.Now(),
		Assistant: true,
	}, nil
}
