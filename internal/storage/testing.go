package storage

import (
	"testing"

	"github.com/jmoiron/sqlx"
)

func TestInMemoDb(t *testing.T) *sqlx.DB {
	t.Helper()

	db := sqlx.MustConnect("sqlite3", ":memory:")

	schema := `
		CREATE TABLE visited (
			id integer PRIMARY KEY,
			title VARCHAR NOT NULL,
            add_date VARCHAR NOT NULL DEFAULT(datetime()),
			poster_url varchar,
			status integer NOT NULL DEFAULT 1
		);
	`

	db.MustExec(schema)

	return db
}
