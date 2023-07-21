package main

import (
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func TestNewLOINCTerminologyService(t *testing.T) {
	var err error

	err = godotenv.Load()

	if err != nil {
		t.Error("error loading .env file")
	}

	username := os.Getenv("LOINC_USERNAME")
	pass := os.Getenv("LOINC_PASSWORD")

	if username == "" || pass == "" {
		t.Error("username or password not set: check .env file")
	}

	l, err := NewLOINCTerminologyService(username, pass)

	if err != nil {
		t.Error("error creating LOINC terminology service")
	}

	if l.Name != "LOINC" {
		t.Error("expected LOINC, got", l.Name)
	}
}

func TestLOINCTerminologyService_Ping(t *testing.T) {}
