package document_service

import (
	"fmt"
	document_dto "licor_model/core/modules/document/dto"
	document_repository "licor_model/core/modules/document/repository"
	ai_service "licor_model/core/modules/ollama/service"
	"licor_model/core/util"
	"licor_model/core/util/executor"
	"licor_model/core/util/timer"
)

type DocumentService struct {
	repo *document_repository.DocumentRepository
}

func NewDocumentService() *DocumentService {
	return &DocumentService{
		repo: document_repository.NewDocumentRepository(executor.NewDBExecutor(nil)),
	}
}

func (s *DocumentService) GetDocumentByContextAndOrigin(context string, origin document_dto.OriginEnum) (document_dto.DocumentDto, error) {
	return s.repo.GetDocumentByContextAndOrigin(context, origin)
}

func (s *DocumentService) GetDocumentsBySimiliarity(content string, origins []document_dto.OriginEnum) (docs []document_dto.DocumentDto, err error) {

	timing := timer.NewTimer()
	timing.Start("generating embedding")

	vetor, err := ai_service.GerarEmbedding(content)
	if err != nil {
		return docs, err
	}

	timing.End("generating embedding")

	vetorFormatted := util.BuilderQueryVetor(vetor)

	timing.Start("getting documents")

	docs, err = s.repo.GetDocumentsBySimiliarity(vetorFormatted, origins)

	timing.End("getting documents")

	return docs, err
}

func (s *DocumentService) UpsertDocument(context string, content string, link string, origin document_dto.OriginEnum) error {

	vetor, err := ai_service.GerarEmbedding(content)

	if err != nil {
		return err
	}

	vetorSQL := util.BuilderQueryVetor(vetor)
	document, err := s.GetDocumentByContextAndOrigin(context, origin)

	if err != nil {
		err = s.repo.CreateDocument(context, content, link, vetorSQL, origin)
		fmt.Println("criando")
	} else {
		document.Link = link
		document.Content = content
		err = s.repo.UpdateDocument(document, vetorSQL)
	}
	return err
}

func (s *DocumentService) DeleteAllFromOrigin(origin document_dto.OriginEnum) error {
	return s.repo.DeleteAllDocumentsFromOrigin(origin)
}
