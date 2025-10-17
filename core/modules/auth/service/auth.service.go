package auth_service

import (
	"errors"
	"fmt"
	auth_dto "licor_model/core/modules/auth/dto"
	auth_repository "licor_model/core/modules/auth/repository"
	"licor_model/core/util/executor"
	guard_util "licor_model/core/util/jwt"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo *auth_repository.AuthRepository
}

func NewAuthService() *AuthService {
	return &AuthService{repo: auth_repository.NewAuthRepository(executor.NewDBExecutor(nil))}
}

func (s *AuthService) GetUserByID(id string) (auth_dto.UserDto, error) {
	user, err := s.repo.GetUserByID(id)
	return user, err
}

func (s *AuthService) Register(registerDto auth_dto.RegisterRequestDto) (auth_dto.AuthResponseDto, error) {
	var response auth_dto.AuthResponseDto

	exists, err := s.repo.EmailExists(registerDto.Email)
	if err != nil {
		return response, err
	}

	if exists {
		return response, errors.New("email ja cadastrado")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registerDto.Password), bcrypt.DefaultCost)
	if err != nil {
		return response, err
	}

	userID, err := s.repo.CreateUser(registerDto, string(hashedPassword))
	if err != nil {
		return response, err
	}

	user, err := s.repo.GetUserByID(userID)
	if err != nil {
		return response, err
	}

	jwtClaims := auth_dto.JWTClaimsDto{
		ID:    user.ID,
		Email: user.Email,
	}

	token, err := guard_util.GenerateJwt(jwtClaims, 1440)
	if err != nil {
		return response, err
	}

	response.Token = token
	response.User = user
	return response, nil
}

func (s *AuthService) Login(loginDto auth_dto.LoginRequestDto) (auth_dto.AuthResponseDto, error) {
	var response auth_dto.AuthResponseDto

	user, hashedPassword, err := s.repo.GetUserByEmail(loginDto.Email)
	if err != nil {
		fmt.Println("user", err.Error())
		return response, errors.New("credenciais inválidas")
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(loginDto.Password))
	if err != nil {
		fmt.Println("aqui")
		return response, errors.New("credenciais inválidas")
	}

	jwtClaims := auth_dto.JWTClaimsDto{
		ID:    user.ID,
		Email: user.Email,
	}

	token, err := guard_util.GenerateJwt(jwtClaims, 1440)
	if err != nil {
		return response, err
	}

	response.Token = token
	response.User = user
	return response, nil
}
