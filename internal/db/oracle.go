package db

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"time"

	_ "github.com/godror/godror"
)

var DB *sql.DB

type DBconfigOracle struct {
	User     string
	Password string
	Host     string
	Port     int
	Service  string
}

func loadConfigDBOracle(environment string) DBconfigOracle {
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port, err := strconv.Atoi(os.Getenv("DB_PORT"))
	service := os.Getenv("DB_SERVICE")
	println("User", user)
	println("PASS", password)
	println("service", service)
	println("HOST", host)
	if err != nil {
		println("Error al cargar la variable: ", err)
	}

	if environment == "PROD" {
		return DBconfigOracle{
			User:     user,
			Password: password,
			Host:     host,
			Port:     port,
			Service:  service,
		}
	} else {
		return DBconfigOracle{
			User:     user,
			Password: password,
			Host:     host,
			Port:     port,
			Service:  service,
		}
	}

}

func ConnectOracle(environment string) (*sql.DB, error) {
	cfg_ocp_erp_qa := loadConfigDBOracle(environment)
	dsn := fmt.Sprintf("%s/%s@%s:%d/%s", cfg_ocp_erp_qa.User, cfg_ocp_erp_qa.Password, cfg_ocp_erp_qa.Host, cfg_ocp_erp_qa.Port, cfg_ocp_erp_qa.Service)
	println("DSN:", dsn)
	db, err := sql.Open("godror", dsn)
	if err != nil {
		return nil, fmt.Errorf("error al abrir conexión: %w", err)
	}

	db.SetMaxOpenConns(50)

	db.SetMaxIdleConns(50)

	db.SetConnMaxLifetime(5 * time.Minute)

	db.SetConnMaxIdleTime(2 * time.Minute)

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("error de ping a Oracle: %w", err)
	}

	fmt.Printf("Pool de conexiones inicializado (Max: %d)\n", 50)

	return db, nil
}
func InitOracleERP(environment string) error {
	var err error
	DB, err = ConnectOracle(environment)
	return err
}
