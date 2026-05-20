package storage

import (
	"testing"

	"github.com/jmoiron/sqlx"
)

func TestInMemoDb(t *testing.T) *sqlx.DB {
	t.Helper()

	db, err := NewDB(":memory:")
	if err != nil {
		t.Errorf("failed to open test db: %s", err)
	}

	return db
}
