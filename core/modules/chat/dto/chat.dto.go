package chat_dto

import "time"

type MensagemDto struct {
	Id         int       `json:"id"`
	IdChat     string    `json:"id_chat"`
	Conteudo   string    `json:"conteudo"`
	CriadoEm   time.Time `json:"criado_em"`
	Assistente bool      `json:"assistente"`
}

type NovaMensagemDto struct {
	Conteudo string `json:"conteudo" validator:"required"`
}

type NovoChatDto struct {
	IdUsuario int    `json:"id_usuario" validator:"required"`
	Titulo    string `json:"titulo" validator:"required"`
}

type ChatDto struct {
	Id           string    `json:"id"`
	IdUsuario    int       `json:"id_usuario"`
	CriadoEm     time.Time `json:"criado_em"`
	AtualizadoEm time.Time `json:"atualizado_em"`
	Titulo       string    `json:"titulo"`
}

type FiltrosDto struct {
	Pagina     int    `query:"pagina" validator:"optional"`
	Quantidade int    `query:"quantidade" validator:"optional"`
	Pesquisa   string `query:"pesquisa" validator:"optional"`
}
