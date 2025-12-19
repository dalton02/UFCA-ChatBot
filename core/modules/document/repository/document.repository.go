package document_repository

import (
	document_dto "licor_model/core/modules/document/dto"
	"licor_model/core/server/shared"
	"licor_model/core/util/executor"

	"github.com/doug-martin/goqu/v9"
	"github.com/lib/pq"
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

func (r *DocumentRepository) DeleteAllDocumentsFromOrigin(origin document_dto.OriginEnum) error {

	query := `DELETE FROM document WHERE origin = $1`
	_, err := r.executor.Exec(query, origin)
	return err

}

func (r *DocumentRepository) GetDocumentByContextAndOrigin(context string, origin document_dto.OriginEnum) (document_dto.DocumentDto, error) {

	var document document_dto.DocumentDto
	query := `SELECT context,content,link,id FROM document WHERE context=$1 AND origin = $2`
	result := r.executor.QueryRow(query, context, origin)

	err := result.Scan(&document.Context, &document.Content, &document.Link, &document.ID)
	if err != nil {
		return document, err
	}
	return document, nil
}

func (r *DocumentRepository) GetDocumentsBySimiliarity(vetor string, origins []document_dto.OriginEnum) (docs []document_dto.DocumentDto, err error) {
	if len(origins) == 0 {
		return docs, nil
	}

	query := `SELECT context,content,link FROM document WHERE origin = ANY($1) ORDER BY embedding <-> ` + vetor + ` LIMIT 8`
	result, err := r.executor.Query(query, pq.Array(origins))
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

func (r *DocumentRepository) CreateDocument(context, content, link, vetorSQL string, origin document_dto.OriginEnum) error {
	id := ksuid.New().String()
	queryInsert := `INSERT INTO document (id,context,content,link,origin,embedding) VALUES ($1,$2,$3,$4,$5,` + vetorSQL + `)`
	_, err := r.executor.Exec(queryInsert, id, context, content, link, origin)
	return err
}

// UpdateDocument atualiza um documento existente pelo ID
func (r *DocumentRepository) UpdateDocument(doc document_dto.DocumentDto, vetorSQL string) error {
	queryUpdate := `UPDATE document SET context = $1, content = $2, link = $3, origin = $4, embedding = ` + vetorSQL + ` WHERE id = $5`
	_, err := r.executor.Exec(queryUpdate, doc.Context, doc.Content, doc.Link, doc.Origin, doc.ID)
	return err
}
