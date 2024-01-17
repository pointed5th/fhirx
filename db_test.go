package fhird

import (
	"testing"
)

func TestNewDbHandler(t *testing.T) {
	db, err := NewDB()

	if err != nil {
		t.Error(err)
	}

	if db == nil {
		t.Error("db is nil")
	}

	err = db.Ping()

	if err != nil {
		t.Error(err)
	}

	err = db.Close()

	if err != nil {
		t.Error(err)
	}
}
