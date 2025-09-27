package auth_dto

import util_dto "licor_model/core/util/dto"

type LoginRequestDto struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type RegisterRequestDto struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type AuthResponseDto struct {
	Token string  `json:"token"`
	User  UserDto `json:"user"`
}

type UserDto struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	util_dto.TimeStampDefaultDB
}

type JWTClaimsDto struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
}
