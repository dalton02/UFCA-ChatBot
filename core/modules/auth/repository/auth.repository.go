package auth_repository

import (
	auth_dto "licor_model/core/modules/auth/dto"
	"licor_model/core/server/shared"
	"licor_model/core/util/executor"

	"github.com/doug-martin/goqu/v9"
	"github.com/segmentio/ksuid"
)

type AuthRepository struct {
	builder  goqu.DialectWrapper // para criar queries
	executor executor.Executor   // p executar queries
}

func NewAuthRepository(exec executor.Executor) *AuthRepository {
	return &AuthRepository{
		builder:  shared.Builder,
		executor: exec,
	}
}

func (r *AuthRepository) CreateUser(user auth_dto.RegisterRequestDto, hashedPassword string) (string, error) {
	id := ksuid.New().String()

	sql, args, err := r.builder.
		Insert("user").
		Cols("id", "name", "login", "password").
		Vals(goqu.Vals{id, user.Name, user.Email, hashedPassword}).
		ToSQL()

	if err != nil {
		return "", err
	}

	_, err = r.executor.Exec(sql, args...)
	return id, err
}

func (r *AuthRepository) GetUserByEmail(email string) (auth_dto.UserDto, string, error) {
	var user auth_dto.UserDto
	var password string

	sql, args, err := r.builder.
		From("user").
		Select("id", "name", "login", "password").
		Where(goqu.Ex{"login": email}).
		ToSQL()

	if err != nil {
		return user, "", err
	}

	row := r.executor.QueryRow(sql, args...)
	err = row.Scan(&user.ID, &user.Name, &user.Email, &password, &user.CreatedAt, &user.UpdatedAt)

	return user, password, err
}

func (r *AuthRepository) GetUserByID(id string) (auth_dto.UserDto, error) {
	var user auth_dto.UserDto

	sql, args, err := r.builder.
		From("users").
		Select("id", "name", "email", "created_at", "updated_at").
		Where(goqu.Ex{"id": id}).
		ToSQL()

	if err != nil {
		return user, err
	}

	row := r.executor.QueryRow(sql, args...)
	err = row.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt)

	return user, err
}

func (r *AuthRepository) EmailExists(email string) (bool, error) {
	var contador int

	sql, args, err := r.builder.
		From("user").
		Select(goqu.COUNT("*")).
		Where(goqu.Ex{"login": email}).
		ToSQL()

	if err != nil {
		return false, err
	}

	row := r.executor.QueryRow(sql, args...)
	err = row.Scan(&contador)

	return contador > 0, err
}
