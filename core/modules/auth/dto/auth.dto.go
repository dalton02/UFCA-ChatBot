package auth_dto

type LoginRequestDto struct {
	Login string `json:"login" validator:"required"`
	Senha string `json:"senha" validator:"required"`
}
