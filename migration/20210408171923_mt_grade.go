package migrations

import (
	"database/sql"
	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(upMtGrade, downMtGrade)
}

func upMtGrade(tx *sql.Tx) error {
	// This code is executed when the migration is applied.
	return nil
}

func downMtGrade(tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	return nil
}
