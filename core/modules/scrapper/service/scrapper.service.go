package scrapper_service

import (
	"fmt"
	document_service "licor_model/core/modules/document/service"
	"strings"

	"github.com/gocolly/colly"
)

type MemoryScrapper struct {
	Materias []string
}

type ScrapperService struct {
	collectorMain     *colly.Collector
	collectorMateria  *colly.Collector
	collectorConteudo *colly.Collector //
	docsService       *document_service.DocumentService
	memory            MemoryScrapper
}

func NewScrapperService(docsService *document_service.DocumentService) *ScrapperService {
	return &ScrapperService{
		collectorMain:     colly.NewCollector(),
		collectorMateria:  colly.NewCollector(),
		collectorConteudo: colly.NewCollector(), //
		docsService:       docsService,
		memory:            MemoryScrapper{},
	}
}

func (s *ScrapperService) Init() {

	fmt.Println("Starting scrapper")
	s.InfoGeralCurso()
	s.InfoCadeirasCurso()
	s.InfoCadeiraDetalhes()

	s.collectorMain.Visit("https://pt.wikiversity.org/wiki/CCT-UFCA/Ci%C3%AAncia_da_Computa%C3%A7%C3%A3o")
	s.collectorMain.Wait()

	materias := strings.Join(s.memory.Materias, " - ")

	s.docsService.UpsertDocument("[CADEIRAS DISPONIVEIS DO CURSO]-[MATERIAS DISPONIVEIS DO CURSO]", materias, "https://pt.wikiversity.org/wiki/CCT-UFCA/Ci%C3%AAncia_da_Computa%C3%A7%C3%A3o")

	fmt.Println("All finished", s.memory.Materias)

}
