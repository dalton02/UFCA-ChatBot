package server

import (
	auth_validator "licor_model/core/modules/auth/dto/validator"
	"licor_model/core/util/interceptor"
)

func RegisterValidations() {

	interceptor.RegisterValidation("strongpassword", auth_validator.ValidateStrongPassword, "deve ser uma senha forte")

}
