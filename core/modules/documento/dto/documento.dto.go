package documento_dto

type DocumentoDto struct {
	Contexto  string      `json:"contexto"`
	Conteudo  string      `json:"conteudo"`
	Link      string      `json:"link"`
	Embedding [][]float32 `json:"embedding"`
	Id        int         `json:"id"`
}
