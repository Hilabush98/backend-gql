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

var DBApp *sql.DB

type DBconfigApp struct {
	User     string
	Password string
	Host     string
	Port     int
	Database string
}

func loadConfigDBApp() DBconfigApp {
	user := os.Getenv("APP_DB_USER")
	password := os.Getenv("APP_DB_PASSWORD")
	host := os.Getenv("APP_DB_HOST")
	port, err := strconv.Atoi(os.Getenv("APP_DB_PORT"))
	database := os.Getenv("APP_DB_NAME")
	if err != nil {
		logs.Error("internal/db/app_service.postgres.go/loadConfigDBApp", fmt.Sprintf("error al leer APP_DB_PORT: %v", err))
	}

		return DBconfigApp{
			User:     user,
			Password: password,
			Host:     host,
			Port:     port,
			Database: database,
		}
}

func ConnectPostgresApp() (*sql.DB, error) {
	cfg := loadConfigDBApp()
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		logs.Error("internal/db/app_service.postgres.go/ConnectPostgresApp", "No se puede abrir la conexion Details: "+fmt.Sprintf("Error: %v", err))
		return nil, err
	}

	db.SetMaxOpenConns(50)
	db.SetMaxIdleConns(50)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(2 * time.Minute)

	if err = db.Ping(); err != nil {
		logs.Error("internal/db/app_service.postgres.go/ConnectPostgresApp", "No se puede hacer ping Details: "+fmt.Sprintf("Error: %v", err))
		return nil, err
	}
	logs.Debug("Pool de conexiones App iniciado", fmt.Sprintf("Pool de conexiones inicializado (Max: %d)", 50))

	return db, nil
}

func InitPostgresApp() error {
	var err error
	DBApp, err = ConnectPostgresApp()
	return err
}
