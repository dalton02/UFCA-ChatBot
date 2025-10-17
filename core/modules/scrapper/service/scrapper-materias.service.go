package scrapper_service

import (
	"fmt"
	ollama_service "licor_model/core/modules/ollama/service"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

func (s *ScrapperService) InfoCadeiraDetalhes() {
	s.collectorMateria.OnHTML("#mw-content-text", func(pagina *colly.HTMLElement) {
		doc := pagina.DOM

		courseName := s.GetCourseName(doc)

		fmt.Println(courseName)

		baseContext := "CADEIRA" + " - " + courseName + " - "
		doc.Find("div.mw-heading.mw-heading2").Each(func(i int, div *goquery.Selection) {

			objective := ""
			syllabus := ""
			h2Objetivos := div.Find("h2#Objetivos")

			h2Ementa := div.Find("h2#Ementa")
			if h2Ementa.Length() > 0 {
				contextName := baseContext + "EMENTA"
				fmt.Println("== EMENTA ==")

				for sib := div.Next(); sib.Length() > 0; sib = sib.Next() {
					if sib.Is("div.mw-heading.mw-heading2") {
						break
					}

					if sib.Is("p") {
						syllabus += sib.Text()
					}

				}

				embedding, err := ollama_service.GerarEmbedding(syllabus)
				if err == nil {
					s.docsService.UpsertDocument(contextName, syllabus, pagina.Request.URL.String(), embedding)
				} else {
					fmt.Println(err.Error())
				}

			}
			if h2Objetivos.Length() > 0 {

				contextName := baseContext + "OBJETIVOS"
				fmt.Println("== OBJETIVOS ==")

				// Itera sobre os irmãos seguintes até o próximo heading
				for sib := div.Next(); sib.Length() > 0; sib = sib.Next() {
					// Parar se acharmos outra div.mw-heading.mw-heading2 (ou outro h2)
					if sib.Is("div.mw-heading.mw-heading2") {
						break
					}

					// Pega todos os <p> que encontrar
					if sib.Is("p") {
						objective += sib.Text()
					}

				}

				embedding, err := ollama_service.GerarEmbedding(objective)
				if err == nil {
					s.docsService.UpsertDocument(contextName, objective, pagina.Request.URL.String(), embedding)
				} else {
					fmt.Println(err.Error())
				}

			}
		})
	})
}

func (s *ScrapperService) GetCourseName(dom *goquery.Selection) (courseName string) {

	dom.Find("tr").EachWithBreak(func(i int, tr *goquery.Selection) bool {
		texto := strings.ToLower(tr.Text())

		// Verifica se a linha contém "componente curricular"
		if strings.Contains(texto, "componente curricular") {
			// Pega a segunda <td> (índice 1)
			td := tr.Find("td").Eq(1)
			courseName = strings.TrimSpace(td.Text())

			fmt.Println("Nome da cadeira:", courseName)
			return false
		}
		return true
	})

	if courseName == "" {
		fmt.Println("Não foi possível encontrar o nome da cadeira.")
	}
	return courseName
}
