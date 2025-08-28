package db

import (
	"database/sql"
	"fmt"

	"github.com/sirupsen/logrus"
)

const (
	DefaultHost               = "127.0.0.1"
	DefaultPort               = "5432"
	DefaultDatabaseName       = "altron_db"
	DefaultUser               = "postgres"
	DefaultSSLMode            = "disable"
	DefaultMigrationDirectory = "migrations"
	DefaultDriverName         = "postgres"
)

type Config struct {
	Host               string `yaml:"db-host"`
	Port               string `yaml:"db-port"`
	DatabaseName       string `yaml:"db-name"`
	User               string `yaml:"db-user"`
	Password           string `yaml:"db-password"`
	Schema             string `yaml:"db-schema"`
	MigrationDirectory string `yaml:"db-migrations-directory"`
	SslMode            string `yaml:"db-ssl"`
	DriverName         string `yaml:"db-driver"`
}

func (db *Config) String() string {
	return fmt.Sprintf("Connecting to DB on %s:%s/%s as '%s' ...", db.Host, db.Port, db.DatabaseName, db.User)
}

// DSN - Data Source Name or connection string
func (db *Config) DSN() string {
	dsn := fmt.Sprintf(
		"host=%s port=%s dbname=%s sslmode=%s user=%s password=%s",
		db.Host, db.Port, db.DatabaseName, db.SslMode, db.User, db.Password,
	)
	if db.Schema != "" {
		dsn += fmt.Sprintf(" search_path=%s,public", db.Schema)
	}
	return dsn
}

func GetConnect(cfg *Config, log *logrus.Logger) (*sql.DB, error) {
	log.Infof("Connecting to DB on %s:%s/%s as '%s' ... ", cfg.Host, cfg.Port, cfg.DatabaseName, cfg.User)
	db, err := sql.Open(cfg.DriverName, cfg.DSN())
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return db, nil
}
