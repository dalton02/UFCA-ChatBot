package chat_repository

import (
	"fmt"
	chat_dto "licor_model/core/modules/chat/dto"
	"licor_model/core/server/shared"
	"licor_model/core/util"
	"licor_model/core/util/executor"

	"github.com/doug-martin/goqu/v9"
	"github.com/lib/pq"
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

func (r *ChatRepository) CreateChat(chat chat_dto.CreateChatDto, userID string) (id string, err error) {
	id = ksuid.New().String()
	sql, args, err := r.builder.
		Insert("chat").
		Cols("id", "user_id", "title").
		Vals(
			goqu.Vals{
				id,
				userID,
				chat.Title,
			}).ToSQL()

	if err != nil {
		return id, err
	}
	_, err = r.executor.Exec(sql, args...)

	return id, err
}

// Listar chats (filtros ainda não implementados)
func (r *ChatRepository) ListChat(filters chat_dto.QueryListChatDto, userID string) (response chat_dto.ListChatDto, err error) {
	// Exemplo simples sem aplicar filtros ainda
	response.Limit = filters.Limit
	response.Page = filters.Page

	build := r.builder.From("chat").
		Select("id", "user_id", "title", "created_at", "updated_at").
		Where(goqu.Ex{
			"user_id": userID,
		}).
		Limit(uint(filters.Limit)).
		Offset(uint(filters.Limit * (filters.Page - 1)))

	response.Total, _ = util.CountSQL(r.executor, "DISTINCT id", build)
	sql, args, _ := build.ToSQL()

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

// Listar chats (filtros ainda não implementados)
func (r *ChatRepository) ListMessages(filters chat_dto.QueryListMessageDto, chatID string) (response chat_dto.ListMessageDto, err error) {
	// Exemplo simples sem aplicar filtros ainda
	response.Limit = filters.Limit
	response.Page = filters.Page

	build := r.builder.From("message").
		Select("id", "content", "created_at", "assistant", "links").
		Where(goqu.Ex{
			"chat_id": chatID,
		}).
		Limit(uint(filters.Limit)).
		Offset(uint(filters.Limit * (filters.Page - 1)))

	response.Total, _ = util.CountSQL(r.executor, "DISTINCT id", build)
	sql, args, _ := build.ToSQL()

	rows, queryErr := r.executor.Query(sql, args...)

	if queryErr != nil {
		return response, queryErr
	}

	defer rows.Close()

	var messages []chat_dto.MessageDto = []chat_dto.MessageDto{}

	for rows.Next() {
		var message chat_dto.MessageDto

		message.ChatID = chatID
		if scanErr := rows.Scan(&message.ID, &message.Content, &message.CreatedAt, &message.Assistant, pq.Array(&message.Links)); scanErr != nil {
			return response, scanErr
		}

		messages = append(messages, message)
	}

	response.Data = messages

	return response, nil
}

func (r *ChatRepository) DeleteAllMessagesChat(chatID string) error {
	sql := `DELETE FROM message WHERE chat_id = $1`

	_, err := r.executor.Exec(sql, chatID)

	return err
}

func (r *ChatRepository) DeleteChat(id string) error {
	sql := `DELETE FROM chat WHERE id = $1`

	_, err := r.executor.Exec(sql, id)

	return err
}

// Salvar mensagem
func (r *ChatRepository) CreateMessage(mensagem string, chatID string, assistant bool, links []string) (id string, err error) {
	id = ksuid.New().String()
	sql, args, _ := r.builder.Insert("message").
		Cols("id", "chat_id", "content", "assistant", "links").
		Vals(goqu.Vals{id, chatID, mensagem, assistant, pq.Array(links)}).
		ToSQL()

	_, err = r.executor.Exec(sql, args...)
	return id, err
}
