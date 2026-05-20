package storage

import (
	"github.com/jmoiron/sqlx"
)

func NewDB(url string) (*sqlx.DB, error) {
	db, err := sqlx.Open("sqlite3", url)
	if err != nil {
		return nil, err
	}

	schema := `
		CREATE TABLE IF NOT EXISTS visited (
			id integer PRIMARY KEY,
			title VARCHAR NOT NULL,
            add_date VARCHAR NOT NULL DEFAULT(datetime()),
			poster_url varchar,
			status integer NOT NULL DEFAULT 1 CHECK ( status IN (1,2))
		);
	`

	db.MustExec(schema)

	return db, nil
}
