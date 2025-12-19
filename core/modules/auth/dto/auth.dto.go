package auth_dto

import util_dto "licor_model/core/util/dto"

type LoginRequestDto struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type RegisterRequestDto struct {
	Name     string   `json:"name" validate:"required"`
	Email    string   `json:"email" validate:"required,email"`
	Role     RoleEnum `json:"role"`
	Password string   `json:"password" validate:"required,strongpassword"`
}

type RoleEnum string

const (
	Student RoleEnum = "student"
	Helper  RoleEnum = "helper"
	Manager RoleEnum = "manager"
)

type AuthResponseDto struct {
	Token string  `json:"token"`
	User  UserDto `json:"user"`
}

type UserDto struct {
	ID    string   `json:"id"`
	Name  string   `json:"name"`
	Email string   `json:"email"`
	Role  RoleEnum `json:"role"`
	util_dto.TimeStampDefaultDB
}

type JWTClaimsDto struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}
