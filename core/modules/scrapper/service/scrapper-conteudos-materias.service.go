package scrapper_service

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

// procura  os  links de conteudo nas páginas das matérias
func (s *ScrapperService) InfoConteudos() {
	s.collectorMateria.OnHTML(".mw-parser-output", func(e *colly.HTMLElement) {

		materiaNome := s.GetCourseName(e.DOM)
		if materiaNome == "" {
			fmt.Println("Não foi possível achar o nome da matéria")
			return
		}

		h2Conteudo := e.DOM.Find("h2#Conteúdo")
		if h2Conteudo.Length() == 0 {
			return
		}

		// acha a <div> que envolveo o <h2> p conseguir achar a <ul> "irmã"
		divHeading := h2Conteudo.Closest("div.mw-heading.mw-heading2")
		if divHeading.Length() == 0 {
			return // estrutura inesperada
		}

		ul := divHeading.NextAllFiltered("ul").First()
		if ul.Length() == 0 {
			return // list not founded
		}

		fmt.Printf("[%s] Encontrada seção 'Conteúdo'. Procurando links...\n", materiaNome)

		// itera sobre os <li> -tópicos- que são filhos diretos da <ul>
		ul.ChildrenFiltered("li").Each(func(i int, li *goquery.Selection) {

			link := li.Find("a").First()
			href, exists := link.Attr("href")

			if !exists || len(href) == 0 || strings.Contains(href, "action=edit") {
				return // pula se não for um link de conteúdo válido
			}

			// pega o texto do link
			topicoTitulo := strings.TrimSpace(link.Text())
			if topicoTitulo == "" {
				return
			}

			// prepara o  novo contexto  p o scrapper ddas páginas de conteúdo
			ctx := colly.NewContext()
			ctx.Put("materiaNome", materiaNome)   // nome da matéria
			ctx.Put("topicoTitulo", topicoTitulo) // nome do conteúdo

			// monta a URL e chama o collectorConteudo
			fullURL := "https://pt.wikiversity.org" + href
			fmt.Printf("Visitando : %s\n", topicoTitulo)
			s.collectorConteudo.Request("GET", fullURL, nil, ctx, nil)
		})
	})
}

// visita  os links de conteudo
func (s *ScrapperService) InfoPaginaDeConteudo() {

	s.collectorConteudo.OnHTML(".mw-parser-output", func(e *colly.HTMLElement) {

		materiaNome := e.Request.Ctx.Get("materiaNome")
		topicoTitulo := e.Request.Ctx.Get("topicoTitulo")
		link := e.Request.URL.String()

		baseDocName := fmt.Sprintf("%s / %s", materiaNome, topicoTitulo)

		fmt.Printf("[%s] Processando conteúdo...\n", baseDocName)

		var currentContent strings.Builder
		// aqui eu estou supondo que a primeira seção é SEMPRE a introdução
		var currentSectionTitle string = "Introdução"

		e.DOM.Find("div.mw-parser-output > *").Each(func(i int, el *goquery.Selection) {

			if el.Is("h2") || el.Is("h3") {
				s.saveSection(baseDocName, currentSectionTitle, currentContent.String(), link)

				currentContent.Reset()
				currentSectionTitle = strings.TrimSpace(el.Find("span.mw-headline").Text())

			} else if el.Is("p, ul, dl, div.mw-highlight, div.exemplo") {
				// ignora a tabela de conteudos (o sumário que aparece)
				if !el.Is("div#toc") {
					currentContent.WriteString(el.Text())
					currentContent.WriteString("\n\n")
				}
			}
		})

		// salva  a  ultima seção --o loop  termina antes de salvar o último bloco--
		s.saveSection(baseDocName, currentSectionTitle, currentContent.String(), link)
	})
}

// pega, formata e salva  o conteúdo
func (s *ScrapperService) saveSection(baseName, sectionName, content, url string) {
	if sectionName == "" || strings.Contains(sectionName, "Sumário") {
		return
	}
	content = strings.TrimSpace(content)
	if content == "" {
		return
	}

	docName := fmt.Sprintf("%s / %s", baseName, sectionName)

	fmt.Printf("Salvando Documento: %s (Tamanho: %d)\n", docName, len(content))
	err := s.docsService.UpsertDocument(docName, content, url)
	if err != nil {
		fmt.Printf("Erro ao salvar documento %s: %s\n", docName, err.Error())
	}
}

/* acho que  o scrapper pode acabar travando
o meu scrapper vai chamar o UpsertDocument muitas vezes em sequência pq cada materia pode ter vários tópicos de conteúdo e cada tópico pode ter várias seções

 daí eu estou  achando que vou acabar criando um gargalo gigantesco caso exista uma matéria com muitos tópicos e seções
*/
