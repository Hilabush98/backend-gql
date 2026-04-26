package db

import (
	"backend-gql/internal/logs"
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

	if err != nil {
		logs.Error("internal/db/oracle.go/loadConfigDBOracle", fmt.Sprintf("error al leer DB_PORT: %v", err))

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
	db, err := sql.Open("godror", dsn)
	if err != nil {
		logs.Error("internal/db/oracle.go/ConnectOracle", "No se puede abrir la conexion Details: "+fmt.Sprintf("Error: %v", err))
		return nil, err
	}

	db.SetMaxOpenConns(50)

	db.SetMaxIdleConns(50)

	db.SetConnMaxLifetime(5 * time.Minute)

	db.SetConnMaxIdleTime(2 * time.Minute)

	if err = db.Ping(); err != nil {
		logs.Error("internal/db/oracle.go/ConnectOracle", "No se puede hacer ping Details: "+fmt.Sprintf("Error: %v", err))
		return nil, err
	}
	logs.Debug("Pool de conexiones iniciado", fmt.Sprintf("Pool de conexiones inicializado (Max: %d)", 50))

	return db, nil
}
func InitOracleERP(environment string) error {
	var err error
	DB, err = ConnectOracle(environment)
	return err
}
