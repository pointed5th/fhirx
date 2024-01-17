package fhird

import (
	"database/sql"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	*sql.DB
	File string
}

func NewDB() (*DB, error) {
	file := "./data/fhird.db"

	if _, err := os.Stat(file); os.IsNotExist(err) {
		f, err := os.Create(file)

		if err != nil {
			return nil, err
		}

		defer f.Close()
	}

	db, err := sql.Open("sqlite3", file)

	if err != nil {
		return nil, err
	}

	return &DB{
		File: file,
		DB:   db,
	}, nil
}
