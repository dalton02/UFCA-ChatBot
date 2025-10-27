package scrapper_service

import (
	document_service "licor_model/core/modules/document/service"

	"github.com/gocolly/colly"
)

type ScrapperService struct {
	collectorMain     *colly.Collector
	collectorMateria  *colly.Collector
	collectorConteudo *colly.Collector //
	docsService       *document_service.DocumentService
}

func NewScrapperService(docsService *document_service.DocumentService) *ScrapperService {
	return &ScrapperService{
		collectorMain:     colly.NewCollector(),
		collectorMateria:  colly.NewCollector(),
		collectorConteudo: colly.NewCollector(), //
		docsService:       docsService,
	}
}

func (s *ScrapperService) Init() {

	s.InfoGeralCurso()
	s.InfoCadeirasCurso()
	s.InfoCadeiraDetalhes()
	s.InfoConteudos()        //
	s.InfoPaginaDeConteudo() //

	s.collectorMain.Visit("https://pt.wikiversity.org/wiki/CCT-UFCA/Ci%C3%AAncia_da_Computa%C3%A7%C3%A3o")
}
