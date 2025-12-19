package document_dto

type DocumentDto struct {
	Context   string      `json:"context"`
	Content   string      `json:"content"`
	Origin    OriginEnum  `json:"origin"`
	Link      string      `json:"link"`
	Embedding [][]float32 `json:"embedding"`
	ID        string      `json:"id"`
}

type OriginEnum string

const (
	Wikiversity  OriginEnum = "wikiversity"
	JsonStudents OriginEnum = "json_students"
	JsonServers  OriginEnum = "json_servers"
)
