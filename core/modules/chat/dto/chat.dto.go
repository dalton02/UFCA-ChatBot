package chat_dto

import (
	util_dto "licor_model/core/util/dto"
	"time"
)

type MessageDto struct {
	ID        int       `json:"id"`
	ChatID    string    `json:"chatId"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createdAt"`
	Assistant bool      `json:"assistant"`
}

type CreateMensagemDto struct {
	Content string `json:"content" validator:"required"`
}

type CreateChatDto struct {
	UserID string `json:"userId" validator:"required"`
	Title  string `json:"title" validator:"required"`
}

type ChatDto struct {
	ID     string `json:"id"`
	UserID string `json:"userId"`
	Title  string `json:"title"`
	util_dto.TimeStampDefaultDB
}

type ListChatDto struct {
	util_dto.QueryPaginationDto
}
