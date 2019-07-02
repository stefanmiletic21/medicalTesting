package db

import (
	"database/sql"
	"medicalTesting/config"

	_ "github.com/lib/pq"
)

var handle *sql.DB

func InitializeDb() error {
	dbHandle, err := sql.Open("postgres", config.GetDatabaseConnectionString())
	if err != nil {
		return err
	}

	dbHandle.SetMaxIdleConns(config.GetDatabaseMaxIdleConnections())
	dbHandle.SetMaxOpenConns(config.GetDatabaseMaxOpenConnections())
	dbHandle.SetConnMaxLifetime(config.GetDatabaseConnectionMaxLifetime())

	err = dbHandle.Ping()
	if err != nil {
		return err
	}

	handle = dbHandle

	return nil
}
