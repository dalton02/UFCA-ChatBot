package chat_dto

import "time"

type MensagemDto struct {
	Id         int       `json:"id"`
	IdChat     string    `json:"id_chat"`
	Conteudo   string    `json:"conteudo"`
	CriadoEm   time.Time `json:"criado_em"`
	Assistente bool      `json:"assistente"`
}

type NovoChatDto struct {
	IdUsuario int    `json:"id_usuario"`
	Titulo    string `json:"titulo"`
}

type ChatDto struct {
	Id           string    `json:"id"`
	IdUsuario    int       `json:"id_usuario"`
	CriadoEm     time.Time `json:"criado_em"`
	AtualizadoEm time.Time `json:"atualizado_em"`
	Titulo       string    `json:"titulo"`
}

type FiltrosDto struct {
	Pagina     int    `json:"pagina"`
	Quantidade int    `json:"quantidade"`
	Pesquisa   string `json:"pesquisa"`
}
