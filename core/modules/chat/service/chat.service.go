package chat_service

import (
	"fmt"
	chat_dto "licor_model/core/modules/chat/dto"
	"licor_model/core/server/shared"
)

func CriarChatDB(chat chat_dto.NovoChatDto) (string, error) {
	var id string
	query := `INSERT INTO chats (id_usuario,titulo) VALUES ($1, $2) RETURNING id`
	result := shared.DB.QueryRow(query, chat.IdUsuario, chat.Titulo)
	err := result.Scan(&id)
	if err != nil {
		return id, err
	}
	return id, nil
}

func ListarChatDB(filtros chat_dto.FiltrosDto) {

}

func BuscarChatDB(id string) (chat_dto.ChatDto, error) {
	var chat chat_dto.ChatDto
	query := `SELECT id,id_usuario,titulo,criado_em,atualizado_em FROM chats WHERE id = $1`
	result, err := shared.DB.Query(query, id)
	if err != nil {
		return chat, err
	}
	if !result.Next() {
		return chat, fmt.Errorf("nenhum chat encontrado")
	}
	err = result.Scan(&chat.Id, &chat.IdUsuario, &chat.Titulo, &chat.CriadoEm, &chat.AtualizadoEm)
	if err != nil {
		return chat, err
	}
	return chat, nil

}

func BuscarMensagensDB(id int) {

}

func SalvarMensagemDB(mensagem string, idChat string, assistente bool) (int, error) {
	var id int
	query := `INSERT INTO mensagens (id_chat,conteudo,assistente) VALUES ($1,$2,$3) RETURNING id`
	result := shared.DB.QueryRow(query, idChat, mensagem, assistente)
	err := result.Scan(&id)
	if err != nil {
		return id, err
	}
	return id, nil
}
