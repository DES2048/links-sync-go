package storage_test

import (
	"links-sync-go/internal/storage"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVisitedRepo_Add(t *testing.T) {
	db := storage.TestInMemoDb(t)

	repo := storage.NewVisitedDbRepo(db)

	data := &storage.DbCreateVisited{
		Title:     "test topic",
		PosterUrl: "some url",
		Status:    1,
	}

	got, err := repo.Add(data)

	assert.Nil(t, err)

	assert.Greater(t, got.Id, 0)
}
