package auth_controller

import (
	"net/http"

	"github.com/dalton02/licor/httpkit"
)

// Login do usuário
// @Summary Autentica um usuário
// @Description Recebe login e senha para autenticar
// @Tags Autenticação
// @Accept json
// @Produce json
// @Param credentials body auth_dto.LoginRequestDto true "User credentials"
// @Router /login [post]
func Login(response http.ResponseWriter, request *http.Request) httpkit.HttpMessage {
	return httpkit.AppSuccess("Login realizado com sucesso", nil)
}

// Cadastro do usuário
// @Summary Cadastra um usuário
// @Description Eu acho que cadastra um usuário
// @Tags Autenticação
// @Accept json
// @Produce json
// @Param credentials body auth_dto.LoginRequestDto true "User credentials"
// @Router /cadastrar [post]
func Cadastro(response http.ResponseWriter, request *http.Request) httpkit.HttpMessage {
	return httpkit.AppSuccess("Cadastro realizado com sucesso", nil)
}
