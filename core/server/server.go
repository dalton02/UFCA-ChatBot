package server

import (
	"fmt"
	ollama_dto "licor_model/core/modules/ollama/dto"
	ollama_service "licor_model/core/modules/ollama/service"
	_ "licor_model/docs"
	"os"

	"github.com/dalton02/licor/licor"
)

func MainServer() {

	message, _ := ollama_service.SendRequest(ollama_dto.Request{
		Model:  os.Getenv("OLLAMA_MODEL"),
		Stream: false,
		Messages: []ollama_dto.Message{
			{Role: "user", Content: "ALOHA"},
		},
	})

	fmt.Println(message.Message)

	licor.Init("4000")

}
