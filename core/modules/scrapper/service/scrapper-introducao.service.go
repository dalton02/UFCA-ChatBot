package scrapper_service

import (
	"fmt"
	ollama_service "licor_model/core/modules/ollama/service"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

func (s *ScrapperService) InfoCadeirasCurso() {
	s.collectorMain.OnHTML(".wikitable", func(table *colly.HTMLElement) {
		tabelaCadeiras := false
		table.ForEach("tr", func(i int, row *colly.HTMLElement) {
			if strings.Contains(row.Text, "Semestre") && i == 0 {
				tabelaCadeiras = true
			}

			if i <= 3 || (i != 0 && !tabelaCadeiras) {
				return
			}

			row.ForEach("td", func(j int, column *colly.HTMLElement) {
				if j != 1 {
					return
				}
				column.ForEach("a", func(i int, link *colly.HTMLElement) {
					linkFormatado := link.DOM.AttrOr("href", "")
					if len(linkFormatado) > 0 {
						s.collectorMateria.Visit(`https://pt.wikiversity.org` + linkFormatado)
					}
				})
			})
		})
	})
}

func (s *ScrapperService) InfoCadeiraDetalhes() {
	s.collectorMateria.OnHTML("#mw-content-text", func(pagina *colly.HTMLElement) {

	})
}

func (s *ScrapperService) InfoGeralCurso() {
	s.collectorMain.OnHTML("#mw-content-text", func(e *colly.HTMLElement) {
		e.ForEach("p", func(i int, h *colly.HTMLElement) {
			if i > 1 {
				return
			}
			embedding, err := ollama_service.GerarEmbedding(h.Text)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			var builder strings.Builder
			builder.WriteString("INTRODUÇÃO CURSO CC")
			builder.WriteString("-PARTE-")
			builder.WriteString(strconv.Itoa(i))

			err = s.docsService.UpsertDocument(builder.String(), h.Text, "Link", embedding)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
		})

	})

}
