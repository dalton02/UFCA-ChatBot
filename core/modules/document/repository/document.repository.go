package document_repository

import (
	document_dto "licor_model/core/modules/document/dto"
	"licor_model/core/server/shared"
	"licor_model/core/util/executor"

	"github.com/doug-martin/goqu/v9"
	"github.com/segmentio/ksuid"
)

type DocumentRepository struct {
	builder  goqu.DialectWrapper
	executor executor.Executor
}

func NewDocumentRepository(exec executor.Executor) *DocumentRepository {
	return &DocumentRepository{
		builder:  shared.Builder,
		executor: exec,
	}
}

func (r *DocumentRepository) GetDocumentByContext(context string) (document_dto.DocumentDto, error) {

	var document document_dto.DocumentDto
	query := `SELECT context,content,link,id FROM document WHERE context=$1`
	result := r.executor.QueryRow(query, context)

	err := result.Scan(&document.Context, &document.Content, &document.Link, &document.ID)
	if err != nil {
		return document, err
	}
	return document, nil
}

func (r *DocumentRepository) GetDocumentsBySimiliarity(vetor string) (docs []document_dto.DocumentDto, err error) {
	query := `SELECT context,content,link FROM document ORDER BY embedding <-> ` + vetor + ` LIMIT 8`
	result, err := r.executor.Query(query)
	if err != nil {
		return docs, err
	}
	for result.Next() {
		var doc document_dto.DocumentDto
		result.Scan(&doc.Context, &doc.Content, &doc.Link)
		docs = append(docs, doc)
	}
	return docs, nil
}

func (r *DocumentRepository) CreateDocument(context string, content string, link string, vetorSQL string) error {
	id := ksuid.New().String()
	queryInsert := `INSERT INTO document (id,context,content,link,embedding) VALUES ($1,$2,$3,$4,` + vetorSQL + `)`
	_, err := r.executor.Exec(queryInsert, id, context, content, link)
	return err
}

// UpdateDocument atualiza um documento existente pelo ID
func (r *DocumentRepository) UpdateDocument(doc document_dto.DocumentDto, vetorSQL string) error {
	queryUpdate := `UPDATE document SET context = $1, content = $2, link = $3, embedding = ` + vetorSQL + ` WHERE id = $4`
	_, err := r.executor.Exec(queryUpdate, doc.Context, doc.Content, doc.Link, doc.ID)
	return err
}
