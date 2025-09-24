package shared

import (
	"database/sql"

	"github.com/doug-martin/goqu/v9"
	"github.com/gin-contrib/cors"
)

var DB *sql.DB

var Builder = goqu.Dialect("postgres")
var Cors = cors.New(cors.Config{
	AllowMethods:     []string{"*"},
	AllowHeaders:     []string{"*"},
	ExposeHeaders:    []string{"*"},
	MaxAge:           1800,
	AllowCredentials: true,
	AllowAllOrigins:  true,
})

func SetDB(db *sql.DB) {
	DB = db
}
func GetDB() *sql.DB {
	return DB
}
