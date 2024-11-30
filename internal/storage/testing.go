package storage

import (
	"testing"

	"github.com/jmoiron/sqlx"
)

func TestInMemoDb(t *testing.T) *sqlx.DB {
	t.Helper()

	db := sqlx.MustConnect("sqlite3", ":memmory")

	schema := `
		CREATE TABLE visited (
			id integer PRIMARY KEY AUTOINCREMENT,
			title VARCHAR NOT NULL,
			poster_url varchar,
			status integer NOT NULL
		);
	`

	db.MustExec(schema)

	return db
}
