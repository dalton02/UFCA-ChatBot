package document_dto

type DocumentDto struct {
	Context   string      `json:"context"`
	Content   string      `json:"content"`
	Link      string      `json:"link"`
	Embedding [][]float32 `json:"embedding"`
	ID        int         `json:"id"`
}
