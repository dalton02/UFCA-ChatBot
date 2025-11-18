package scrapper_service

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

func (s *ScrapperService) InfoCadeiraDetalhes() {
	s.collectorMateria.OnHTML("#mw-content-text", func(pagina *colly.HTMLElement) {
		doc := pagina.DOM

		courseName := s.GetCourseName(doc)

		baseContext := "CADEIRA" + " - " + courseName + " - "

		fmt.Println(courseName, " -  VISITADO")

		doc.Find("div.mw-heading.mw-heading2").Each(func(i int, div *goquery.Selection) {

			h2Objetivos := div.Find("h2#Objetivos")
			h2Avalicao := div.Find("h2#Avaliação")
			h2Ementa := div.Find("h2#Ementa")
			h2Metodologia := div.Find("h2#Metodologia")
			if h2Ementa.Length() > 0 {
				contextName := baseContext + "EMENTA"
				syllabus := s.GetParagraphsAfter(div)
				s.docsService.UpsertDocument(contextName, syllabus, pagina.Request.URL.String())
			}
			if h2Objetivos.Length() > 0 {
				contextName := baseContext + "OBJETIVOS"
				objective := s.GetParagraphsAfter(div)
				s.docsService.UpsertDocument(contextName, objective, pagina.Request.URL.String())
			}

			if h2Avalicao.Length() > 0 {
				contextName := baseContext + "AVALIACAO"
				content := s.GetParagraphsAfter(div)
				s.docsService.UpsertDocument(contextName, content, pagina.Request.URL.String())
			}
			if h2Metodologia.Length() > 0 {
				contextName := baseContext + "METODOLOGIA"
				content := s.GetParagraphsAfter(div)
				s.docsService.UpsertDocument(contextName, content, pagina.Request.URL.String())
			}
		})
	})
}

func (s *ScrapperService) GetParagraphsAfter(div *goquery.Selection) (content string) {
	for sib := div.Next(); sib.Length() > 0; sib = sib.Next() {
		// Parar se acharmos outra div.mw-heading.mw-heading2 (ou outro h2)
		if sib.Is("div.mw-heading.mw-heading2") {
			break
		}

		// Pega todos os <p> que encontrar
		if sib.Is("p") {
			content += sib.Text()
		}

	}
	return content
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
