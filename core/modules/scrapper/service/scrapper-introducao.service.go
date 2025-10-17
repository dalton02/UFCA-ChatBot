package scrapper_service

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

const LIMIT_VISITS = 5

func (s *ScrapperService) InfoCadeirasCurso() {
	visitsMade := 0 //impedir meu pczin de ser bloqueado

	s.collectorMain.OnHTML(".wikitable", func(table *colly.HTMLElement) {
		tabelaCadeiras := false

		table.ForEach("tr", func(i int, row *colly.HTMLElement) {
			if strings.Contains(row.Text, "Semestre") && i == 0 {
				tabelaCadeiras = true
			}

			if i <= 2 || (i != 0 && !tabelaCadeiras) {
				return
			}

			row.ForEach("td", func(j int, column *colly.HTMLElement) {
				if j != 1 {
					return
				}
				column.ForEach("a", func(i int, link *colly.HTMLElement) {
					linkFormatado := link.DOM.AttrOr("href", "")
					if len(linkFormatado) > 0 {
						fmt.Println(linkFormatado)
						if visitsMade < LIMIT_VISITS {
							s.collectorMateria.Visit(`https://pt.wikiversity.org` + linkFormatado)
						}
						visitsMade += 1
					}
				})
			})
		})
	})
}

func (s *ScrapperService) InfoGeralCurso() {
	s.collectorMain.OnHTML("#mw-content-text", func(e *colly.HTMLElement) {
		e.ForEach("p", func(i int, h *colly.HTMLElement) {
			if i > 1 {
				return
			}
			var builder strings.Builder
			builder.WriteString("INTRODUÇÃO CURSO CC")
			builder.WriteString("-PARTE-")
			builder.WriteString(strconv.Itoa(i))

			err := s.docsService.UpsertDocument(builder.String(), h.Text, e.Request.URL.String())
			if err != nil {
				fmt.Println(err.Error())

			}
		})

	})

}
