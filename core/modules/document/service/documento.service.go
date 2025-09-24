package document_service

import (
	document_dto "licor_model/core/modules/document/dto"
	document_repository "licor_model/core/modules/document/repository"
	ollama_service "licor_model/core/modules/ollama/service"
	"licor_model/core/util"
	"licor_model/core/util/executor"
)

type DocumentService struct {
	repo *document_repository.DocumentRepository
}

func NewDocumentService() *DocumentService {
	return &DocumentService{
		repo: document_repository.NewDocumentRepository(executor.NewDBExecutor(nil)),
	}
}

func (s *DocumentService) GetDocumentByContext(context string) (document_dto.DocumentDto, error) {
	return s.repo.GetDocumentByContext(context)
}

func (s *DocumentService) GetDocumentsBySimiliarity(content string) (docs []document_dto.DocumentDto, err error) {
	vetor, err := ollama_service.GerarEmbedding(content)
	if err != nil {
		return docs, err
	}
	vetorFormatted := util.BuilderQueryVetor(vetor)
	return s.repo.GetDocumentsBySimiliarity(vetorFormatted)
}

func (s *DocumentService) UpsertDocument(context string, content string, link string, vetor [][]float32) error {
	vetorSQL := util.BuilderQueryVetor(vetor)

	document, err := s.GetDocumentByContext(context)

	if err != nil {
		err = s.repo.CreateDocument(context, content, link, vetorSQL)
	} else {
		err = s.repo.UpdateDocument(document, vetorSQL)
	}
	return err
}
