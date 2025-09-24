package chat_repository

import (
	"fmt"
	chat_dto "licor_model/core/modules/chat/dto"
	"licor_model/core/server/shared"
	"licor_model/core/util/executor"

	"github.com/doug-martin/goqu/v9"
	"github.com/segmentio/ksuid"
)

type ChatRepository struct {
	builder  goqu.DialectWrapper
	executor executor.Executor
}

func NewChatRepository(exec executor.Executor) *ChatRepository {
	return &ChatRepository{
		builder:  shared.Builder,
		executor: exec,
	}
}

func (r *ChatRepository) CreateChat(chat chat_dto.CreateChatDto) (id string, err error) {
	id = ksuid.New().String()
	sql, args, err := r.builder.
		Insert("chat").
		Cols("id", "user_id", "title").
		Vals(
			goqu.Vals{
				id,
				chat.UserID,
				chat.Title,
			}).ToSQL()

	if err != nil {
		return id, err
	}
	_, err = r.executor.Exec(sql, args...)

	return id, err
}

// Listar chats (filtros ainda não implementados)
func (r *ChatRepository) ListChat(filtros chat_dto.QueryListChatDto) (response chat_dto.ListChatDto, err error) {
	// Exemplo simples sem aplicar filtros ainda
	sql, args, _ := r.builder.From("chat").
		Select("id", "user_id", "title", "created_at", "updated_at").
		ToSQL()

	rows, queryErr := r.executor.Query(sql, args...)
	if queryErr != nil {
		return response, queryErr
	}
	defer rows.Close()

	var chats []chat_dto.ChatDto = []chat_dto.ChatDto{}
	for rows.Next() {
		var chat chat_dto.ChatDto
		if scanErr := rows.Scan(&chat.ID, &chat.UserID, &chat.Title, &chat.CreatedAt, &chat.UpdatedAt); scanErr != nil {
			return response, scanErr
		}
		chats = append(chats, chat)
	}

	response.Data = chats
	response.Page = filtros.Page
	response.Limit = filtros.Limit

	return response, nil
}

// Buscar chat por ID
func (r *ChatRepository) GetChatByID(id string) (chat_dto.ChatDto, error) {
	sql, args, err := r.builder.From("chat").
		Select("id", "user_id", "title", "created_at", "updated_at").
		Where(goqu.Ex{"id": id}).ToSQL()
	if err != nil {
		return chat_dto.ChatDto{}, err
	}

	rows, queryErr := r.executor.Query(sql, args...)
	if queryErr != nil {
		return chat_dto.ChatDto{}, queryErr
	}
	defer rows.Close()

	if !rows.Next() {
		return chat_dto.ChatDto{}, fmt.Errorf("nenhum chat encontrado")
	}

	var chat chat_dto.ChatDto
	if scanErr := rows.Scan(&chat.ID, &chat.UserID, &chat.Title, &chat.CreatedAt, &chat.UpdatedAt); scanErr != nil {
		return chat_dto.ChatDto{}, scanErr
	}
	return chat, nil
}

// Buscar mensagens de um chat (placeholder)
func (r *ChatRepository) GetMessages(chatID int) ([]string, error) {
	// Ainda não implementado
	return nil, nil
}

// Salvar mensagem
func (r *ChatRepository) CreateMessage(mensagem string, chatID string, assistant bool) (id string, err error) {
	id = ksuid.New().String()
	sql, args, _ := r.builder.Insert("message").
		Cols("id", "chat_id", "content", "assistant").
		Vals(goqu.Vals{id, chatID, mensagem, assistant}).
		ToSQL()

	fmt.Println(sql)
	_, err = r.executor.Exec(sql, args...)
	fmt.Println(err)
	return id, err
}
