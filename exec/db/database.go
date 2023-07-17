//	authd.go - Kahlehla project
//	Copyright 2023 by Kurt Duncan
//	Database handler

package db

// intended for the go exec

import (
	"database/sql"
	"fmt"
)

import (
	"github.com/lib/pq"
	"log"

	_ "github.com/lib/pq"
)

const defaultUserName = "root"
const defaultHostName = "localhost"
const defaultPortNumber = 26257
const defaultDatabaseName = "defaultdb"

type Configuration struct {
	PostgresUserName   string
	PostgresHostName   string
	PostgresPortNumber int
	DatabaseName       string
}

type Database struct {
	SQL *sql.DB
}

func (db *Database) Close() error {
	return db.SQL.Close()
}

func New(config *Configuration) (*Database, error) {
	dataSource := fmt.Sprintf("postgresql://%s@%s:%d/%s?sslmode=disable",
		config.PostgresUserName,
		config.PostgresHostName,
		config.PostgresPortNumber,
		config.DatabaseName)

	log.Printf("Connecting to:%s\n", dataSource)
	db, err := sql.Open("postgres", dataSource)
	if err != nil {
		log.Printf("open failed:%s\n", err)
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		log.Printf("ping failed:%s\n", err)
		return nil, err
	}

	return &Database{
		SQL: db,
	}, nil
}

func LoadConfiguration() (*Configuration, error) {
	//	TODO later on we need to make this configurable
	cfg := Configuration{
		PostgresUserName:   defaultUserName,
		PostgresHostName:   defaultHostName,
		PostgresPortNumber: defaultPortNumber,
		DatabaseName:       defaultDatabaseName,
	}

	return &cfg, nil
}

func isError(err error, msg string) bool {
	if pe, ok := err.(*pq.Error); ok {
		if pe.Code.Name() == msg {
			return true
		}
	}
	return false
}
