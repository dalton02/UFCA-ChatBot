package chat_service

import (
	"fmt"
	chat_dto "licor_model/core/modules/chat/dto"
	chat_repository "licor_model/core/modules/chat/repository"
	document_service "licor_model/core/modules/document/service"
	app "licor_model/core/util/errors"
	"licor_model/core/util/executor"
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

func (s *ChatService) ListChat(filters chat_dto.QueryListChatDto, userID string) (chat_dto.ListChatDto, error) {
	return s.repo.ListChat(filters, userID)
}

func (s *ChatService) DeleteChat(id string, userId string) error {

	chat, err := s.GetChatByID(id)
	if err != nil {
		fmt.Println(err)
		return app.NotFound("chat não encontrado")

	}

	if chat.UserID != userId {
		return app.Forbidden("Nem é teu chat pra deletar amigão")
	}

	err = s.repo.DeleteAllMessagesChat(id)

	if err != nil {
		fmt.Println(err)
		return app.NotFound("chat não encontrado")

	}
	return s.repo.DeleteChat(id)

}

func (s *ChatService) GetChatByID(id string) (chat_dto.ChatDto, error) {
	return s.repo.GetChatByID(id)
}
