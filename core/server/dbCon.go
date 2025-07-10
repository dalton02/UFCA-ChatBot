package server

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
)

func InitConnection() (*sql.DB, error) {
	err := godotenv.Load()
	if err != nil {
		panic(err.Error())
	}

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	num, err := strconv.Atoi(port)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, num, user, password, dbname)

	db, err := sql.Open("postgres", connStr)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		fmt.Println(err)
		db.Close()
		return nil, err
	}

	// Criar string de conexão para golang-migrate
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to initialize driver: %w", err)
	}

	migrationPath, _ := filepath.Abs("migrations")
	_, err = os.Stat(migrationPath)
	if os.IsNotExist(err) {
		fmt.Printf("Diretório %s não existe\n", migrationPath)
		return nil, fmt.Errorf(err.Error())
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://"+migrationPath,
		"postgres", driver,
	)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to initialize migration: %w", err)
	}

	// Aplicar migrações
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		db.Close()
		return nil, fmt.Errorf("failed to apply migrations: %w", err)
	}

	fmt.Println("Migrações aplicadas com sucesso!")
	return db, nil

}
