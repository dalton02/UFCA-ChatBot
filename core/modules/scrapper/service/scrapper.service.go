package scrapper_service

import (
	"fmt"
	documento_service "licor_model/core/modules/documento/service"
	ollama_service "licor_model/core/modules/ollama/service"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

var collectorMain *colly.Collector
var collectorMateria *colly.Collector

func Init() {

	collectorMain = colly.NewCollector()
	collectorMateria = colly.NewCollector()

	infoGeralCurso()
	infoCadeirasCurso()
	infoCadeiraDetalhes()

	collectorMain.Visit("https://pt.wikiversity.org/wiki/CCT-UFCA/Ci%C3%AAncia_da_Computa%C3%A7%C3%A3o")
}

func infoCadeirasCurso() {
	collectorMain.OnHTML(".wikitable", func(table *colly.HTMLElement) {
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
						collectorMateria.Visit(`https://pt.wikiversity.org` + linkFormatado)
					}
				})
			})
		})
	})
}

func infoCadeiraDetalhes() {
	collectorMateria.OnHTML("#mw-content-text", func(pagina *colly.HTMLElement) {

	})
}

func infoGeralCurso() {
	collectorMain.OnHTML("#mw-content-text", func(e *colly.HTMLElement) {
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
			err = documento_service.SalvarEmbeddingDb(builder.String(), h.Text, "Link", embedding)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
		})

	})

}
