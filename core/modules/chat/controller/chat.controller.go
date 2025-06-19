package chat_controller

import (
	"encoding/json"
	"fmt"
	chat_dto "licor_model/core/modules/chat/dto"
	chat_service "licor_model/core/modules/chat/service"
	"net/http"
	"time"

	"github.com/dalton02/licor/httpkit"
)

func NovoChat(response http.ResponseWriter, request *http.Request) httpkit.HttpMessage {

	var chatInfo chat_dto.NovoChatDto
	decoder := json.NewDecoder(request.Body)
	err := decoder.Decode(&chatInfo)
	if err != nil {
		return httpkit.AppBadRequest(err.Error())
	}
	id, err := chat_service.CriarChatDB(chatInfo)
	if err != nil {
		return httpkit.AppBadRequest(err.Error())
	}
	return httpkit.AppSuccess("Chat criado com sucesso", chat_dto.ChatDto{
		Id:           id,
		IdUsuario:    chatInfo.IdUsuario,
		CriadoEm:     time.Now(),
		AtualizadoEm: time.Now(),
		Titulo:       chatInfo.Titulo,
	})

}

func EnviarMensagem(response http.ResponseWriter, request *http.Request) httpkit.HttpMessage {

	return httpkit.AppSuccess("Tudo certo aqui!", make(map[string]interface{}))

}

func BuscarMensagens(response http.ResponseWriter, request *http.Request) httpkit.HttpMessage {

	return httpkit.AppSuccess("Tudo certo aqui!", make(map[string]interface{}))

}

func BuscarChat(response http.ResponseWriter, request *http.Request) httpkit.HttpMessage {

	params, err := httpkit.GetUrlParams(request)
	if err != nil {
		return httpkit.AppBadRequest(err.Error())
	}
	idChat := params.Param["idChat"]
	fmt.Println(idChat)
	result, err := chat_service.BuscarChatDB(idChat)
	if err != nil {
		return httpkit.AppBadRequest(err.Error())
	}
	return httpkit.AppSuccess("Chat criado com sucesso", result)

}
