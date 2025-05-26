package server

import (
	auth_module "licor_model/core/modules/auth"
	_ "licor_model/docs"

	"github.com/dalton02/licor/licor"
)

func MainServer() {
	licor.Public[any, any]("/login").Post(auth_module.Login)
	licor.Public[any, any]("/cadastrar").Post(auth_module.Cadastro)
	licor.Init("4000")
}
