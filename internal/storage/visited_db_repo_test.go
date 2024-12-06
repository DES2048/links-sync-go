package storage_test

import (
	"testing"

	"links-sync-go/internal/storage"

	"github.com/stretchr/testify/assert"
)

func TestVisitedRepo_Add(t *testing.T) {
	db := storage.TestInMemoDb(t)
	defer db.Close()

	repo := storage.NewVisitedDbRepo(db)

	data := &storage.DbCreateVisited{
		Id:        1,
		Title:     "test topic",
		PosterUrl: "some url",
		Status:    1,
	}

	got, err := repo.Add(data)

	assert.Nil(t, err)

	assert.Equal(t, got.Id, data.Id)
	assert.NotEmpty(t, got.AddDate)
}

func TestVisitedRepo_AddBatch(t *testing.T) {
	db := storage.TestInMemoDb(t)
	defer db.Close()

	repo := storage.NewVisitedDbRepo(db)

	data := []*storage.DbCreateVisited{
		{
			Id:        1,
			Title:     "t1",
			PosterUrl: "url1",
			Status:    1,
		},
		{
			Id:        2,
			Title:     "t2",
			PosterUrl: "url2",
			Status:    2,
		},
	}

	expected := make([]*storage.DbVisited, 0, len(data))
	for _, v := range data {
		expected = append(expected, &storage.DbVisited{
			Id:        v.Id,
			Title:     v.Title,
			PosterUrl: v.PosterUrl,
			Status:    v.Status,
		})
	}

	err := repo.AddBatch(data)
	assert.Nil(t, err)

	got, err := repo.List()

	for _, gotElem := range got {
		gotElem.AddDate = ""
	}

	assert.Nil(t, err)

	assert.ElementsMatch(t, expected, got)
}

func TestVisitedRepo_List(t *testing.T) {
	db := storage.TestInMemoDb(t)
	defer db.Close()

	repo := storage.NewVisitedDbRepo(db)

	data := []*storage.DbCreateVisited{
		{
			Id:        1,
			Title:     "t1",
			PosterUrl: "url1",
			Status:    1,
		},
		{
			Id:        2,
			Title:     "t2",
			PosterUrl: "url2",
			Status:    2,
		},
	}

	expected := make([]*storage.DbVisited, 0, len(data))
	for _, v := range data {
		expected = append(expected, &storage.DbVisited{
			Id:        v.Id,
			Title:     v.Title,
			PosterUrl: v.PosterUrl,
			Status:    v.Status,
		})

		_, err := repo.Add(v)
		assert.Nil(t, err)
	}

	got, err := repo.List()

	for _, gotElem := range got {
		gotElem.AddDate = ""
	}

	assert.Nil(t, err)

	assert.ElementsMatch(t, expected, got)
}

func TestVisitedRepo_Get(t *testing.T) {
	db := storage.TestInMemoDb(t)
	defer db.Close()

	repo := storage.NewVisitedDbRepo(db)

	got, err := repo.Get(1)

	assert.ErrorIs(t, err, storage.ErrNotFound)

	expect := &storage.DbCreateVisited{
		Id:    1,
		Title: "title",
	}

	got, err = repo.Add(expect)

	assert.Nil(t, err)
	assert.Equal(t, expect.Id, got.Id)
}

func TestVisitedRepo_PartialUpdate(t *testing.T) {
	db := storage.TestInMemoDb(t)
	defer db.Close()

	repo := storage.NewVisitedDbRepo(db)

	// case ok

	forUpdate := &storage.DbCreateVisited{
		Id:        1,
		Title:     "title",
		PosterUrl: "url1",
		Status:    1,
	}

	_, err := repo.Add(forUpdate)
	assert.Nil(t, err)

	newTitle := "new title"
	updateData := &storage.DbPatchVisited{
		Title: &newTitle,
	}

	err = repo.UpdatePartial(forUpdate.Id, updateData)

	assert.Nil(t, err)
	got, err := repo.Get(forUpdate.Id)
	assert.Nil(t, err)
	assert.Equal(t, got.Title, newTitle)
	// case not found
	// case no fields specified
	// case id in the fields
	// case auto generated field value in fields
	// case invalid field name
	// case invalid field value
}

func TestVisitedRepo_Delete(t *testing.T) {
	db := storage.TestInMemoDb(t)
	defer db.Close()

	repo := storage.NewVisitedDbRepo(db)

	err := repo.Delete(1)

	assert.ErrorIs(t, err, storage.ErrNotFound)

	forDelete := &storage.DbCreateVisited{
		Id:    1,
		Title: "title",
	}

	_, err = repo.Add(forDelete)

	assert.Nil(t, err)

	err = repo.Delete(forDelete.Id)
	assert.Nil(t, err)

	_, err = repo.Get(forDelete.Id)
	assert.ErrorIs(t, err, storage.ErrNotFound)
}
