package db

import (
	"backend-gql/internal/logs"
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"time"

	_ "github.com/lib/pq"
)

var DBAuth *sql.DB

type DBconfigAuth struct {
	User     string
	Password string
	Host     string
	Port     int
	Database string
}

func loadConfigDBAuth() DBconfigAuth {
	user := os.Getenv("AUTH_DB_USER")
	password := os.Getenv("AUTH_DB_PASSWORD")
	host := os.Getenv("AUTH_DB_HOST")
	port, err := strconv.Atoi(os.Getenv("AUTH_DB_PORT"))
	database := os.Getenv("AUTH_DB_NAME")
	fmt.Sprintf(user)

	if err != nil {
		logs.Error("internal/db/auth_service.postgre.go/loadConfigDBAuth", fmt.Sprintf("error al leer AUTH_DB_PORT: %v", err))
	}

			return DBconfigAuth{
			User:     user,
			Password: password,
			Host:     host,
			Port:     port,
			Database: database,
		}
}

func ConnectPostgresAuth() (*sql.DB, error) {
	cfg := loadConfigDBAuth()
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		logs.Error("internal/db/auth_service.postgre.go/ConnectPostgresAuth", "No se puede abrir la conexion Details: "+fmt.Sprintf("Error: %v", err))
		return nil, err
	}

	db.SetMaxOpenConns(50)
	db.SetMaxIdleConns(50)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(2 * time.Minute)

	if err = db.Ping(); err != nil {
		logs.Error("internal/db/auth_service.postgre.go/ConnectPostgresAuth", "No se puede hacer ping Details: "+fmt.Sprintf("Error: %v", err))
		return nil, err
	}
	logs.Debug("Pool de conexiones Auth iniciado", fmt.Sprintf("Pool de conexiones inicializado (Max: %d)", 50))

	return db, nil
}

func InitPostgresAuth() error {
	var err error
	DBAuth, err = ConnectPostgresAuth()
	return err
}
