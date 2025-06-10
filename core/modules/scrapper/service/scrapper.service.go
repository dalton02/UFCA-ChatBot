package scrapper_service

import (
	"fmt"
	ollama_dto "licor_model/core/modules/ollama/dto"
	ollama_service "licor_model/core/modules/ollama/service"
	"os"
	"strings"

	"github.com/gocolly/colly"
)

func Init() {

	c := colly.NewCollector()

	// Find and visit all links

	c.OnHTML("#bodyContent", func(e *colly.HTMLElement) {
		content, _ := e.DOM.Html()
		var contextBuilder strings.Builder
		contextBuilder.WriteString(
			`Responda exclusivamente em português brasileiro. A partir do conteúdo fornecido a seguir, extraia e resuma as seguintes informações: 
		1. Matriz curricular do curso;
		2. Atividades externas disponíveis;
		3. Informações sobre a coordenação do curso.
		Retorne as informações de forma clara, organizada e concisa, utilizando listas ou tópicos quando apropriado. Considere apenas o conteúdo fornecido: 
		---`)
		contextBuilder.WriteString(content)
		contextBuilder.WriteString("---")
		message, err := ollama_service.SendRequest(ollama_dto.Request{
			Model: os.Getenv("OLLAMA_MODEL"),
			Messages: []ollama_dto.Message{
				{Role: "user",
					Content: contextBuilder.String()},
			},
		})
		fmt.Println(message.Message, err)

	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.Visit("https://pt.wikiversity.org/wiki/CCT-UFCA/Ci%C3%AAncia_da_Computa%C3%A7%C3%A3o")
}
