package fhird

import sqlmigrate "github.com/rubenv/sql-migrate"

var MigrationsDir = "./data/migrations"

type Migrator struct {
	Source sqlmigrate.FileMigrationSource
}

func DefaultMigrator() *Migrator {
	return &Migrator{
		Source: sqlmigrate.FileMigrationSource{
			Dir: MigrationsDir,
		},
	}
}
