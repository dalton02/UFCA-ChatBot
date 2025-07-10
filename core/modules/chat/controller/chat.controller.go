package chat_controller

import (
	"encoding/json"
	"fmt"
	chat_dto "licor_model/core/modules/chat/dto"
	chat_service "licor_model/core/modules/chat/service"
	documento_service "licor_model/core/modules/documento/service"
	ollama_dto "licor_model/core/modules/ollama/dto"
	ollama_service "licor_model/core/modules/ollama/service"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dalton02/licor/httpkit"
)

func NovoChat(response http.ResponseWriter, request *http.Request) (httpkit.HttpMessage, bool) {

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

func EnviarMensagem(response http.ResponseWriter, request *http.Request) (httpkit.HttpMessage, bool) {
	params, err := httpkit.GetUrlParams(request)
	if err != nil {
		return httpkit.AppBadRequest(err.Error())
	}
	idChat := params.Param["idChat"]

	var mensagem chat_dto.NovaMensagemDto
	decoder := json.NewDecoder(request.Body)
	decoder.Decode(&mensagem)

	documentos, err := documento_service.BuscarDocumentoPorSimilaridadeDB(mensagem.Conteudo)
	if err != nil {
		return httpkit.AppBadRequest(err.Error())
	}

	var stringBuilder strings.Builder
	stringBuilder.WriteString("REGRAS DO PROMPT:")
	stringBuilder.WriteString("\n1 - Nunca envie links para websites em sua resposta ")
	stringBuilder.WriteString("\n2 - Sempre mantenha o foco do assunto a respeito da universidade e alunos")
	stringBuilder.WriteString("\n3 - Responda a pergunta baseando-se nos seguintes dados: ")
	stringBuilder.WriteString("\n----------INICIO DOS DADOS-----------")
	for _, doc := range documentos {
		stringBuilder.WriteString("\n-----------------------------")
		stringBuilder.WriteString(doc.Conteudo)
		stringBuilder.WriteString("\n-----------------------------")
	}
	stringBuilder.WriteString("\n------ FIM DOS DADOS-------")
	stringBuilder.WriteString("\nPergunta do usuário: ")
	stringBuilder.WriteString(mensagem.Conteudo)
	respostaIA, err := ollama_service.SendRequest(ollama_dto.RequestChatAI{
		Model: os.Getenv("OLLAMA_MODEL"),
		Messages: []ollama_dto.MessageChatAI{
			{
				Content: stringBuilder.String(),
				Role:    "user",
			},
		},
	})
	if err != nil {
		return httpkit.AppBadRequest(err.Error())
	}

	// idMensagem, err := chat_service.SalvarMensagemDB(mensagem.Conteudo, idChat, false)
	// if err != nil {
	// 	if strings.Contains(err.Error(), "foreign key constraint") {
	// 		return httpkit.AppNotFound("Chat não encontrado")
	// 	}
	// 	return httpkit.AppBadRequest(err.Error())
	// }

	return httpkit.AppSuccess("Operação realizada com sucesso", chat_dto.MensagemDto{
		Id:         999,
		IdChat:     idChat,
		Conteudo:   respostaIA.Message.Content,
		CriadoEm:   time.Now(),
		Assistente: true,
	})

}

func BuscarMensagens(response http.ResponseWriter, request *http.Request) (httpkit.HttpMessage, bool) {

	return httpkit.AppSuccess("Tudo certo aqui!", make(map[string]interface{}))

}

func BuscarChat(response http.ResponseWriter, request *http.Request) (httpkit.HttpMessage, bool) {

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
