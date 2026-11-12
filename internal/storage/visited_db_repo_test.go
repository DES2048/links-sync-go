package storage_test

import (
	"strconv"
	"testing"

	"links-sync-go/internal/storage"

	"github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func clearAddDate(data []*storage.DbVisited) {
	// unset AddDate for comparsion
	for _, item := range data {
		item.AddDate = ""
	}
}

func TestVisitedRepo_Add(t *testing.T) {
	db := storage.TestInMemoDb(t)
	defer db.Close()

	repo := storage.NewVisitedDbRepo(db)

	data := &storage.DbCreateVisited{
		Id:        1,
		Title:     "test topic",
		PosterUrl: new("some url"),
		Status:    1,
	}

	got, err := repo.Add(data)

	assert.Nil(t, err)

	assert.Equal(t, got.Id, data.Id)
	assert.Equal(t, got.Title, data.Title)
	assert.Equal(t, got.PosterUrl, data.PosterUrl)
	assert.NotEmpty(t, got.AddDate)
}

func TestVisitedRepo_Add_WithInvalidStatus(t *testing.T) {
	db := storage.TestInMemoDb(t)
	defer db.Close()

	repo := storage.NewVisitedDbRepo(db)

	_, err := repo.Add(&storage.DbCreateVisited{
		Id:     1,
		Title:  "some title",
		Status: 4,
	})

	var sqlErr sqlite3.Error
	// check that is a sqlite3.Error
	assert.ErrorAs(t, err, &sqlErr)
	// check that code is constraint failed
	assert.EqualValues(t, 19, sqlErr.Code)
	// check that extended code is check constraint failed
	assert.EqualValues(t, sqlErr.ExtendedCode, 275)
	// t.Logf("err type %T, value: %#v", err, err)
}

func TestVisitedRepo_Add_WithIDAlreadyExists(t *testing.T) {
	db := storage.TestInMemoDb(t)
	defer db.Close()

	repo := storage.NewVisitedDbRepo(db)

	repo.Add(&storage.DbCreateVisited{
		Id:     1,
		Title:  "some title",
		Status: 1,
	})

	// add another topic with same id
	_, err := repo.Add(&storage.DbCreateVisited{
		Id:     1,
		Title:  "some title2",
		Status: 2,
	})

	var tErr storage.ErrAlreadyExists
	assert.ErrorAs(t, err, &tErr)
}

func TestVisitedRepo_Add_WithDefaultStatus(t *testing.T) {
	db := storage.TestInMemoDb(t)
	defer db.Close()

	repo := storage.NewVisitedDbRepo(db)

	data := &storage.DbCreateVisited{
		Id:        1,
		Title:     "test topic",
		PosterUrl: new("some url"),
		Status:    0,
	}

	got, err := repo.Add(data)

	assert.Nil(t, err)
	assert.Equal(t, 1, got.Status)
}

func TestVisitedRepo_AddBatch(t *testing.T) {
	db := storage.TestInMemoDb(t)
	defer db.Close()

	repo := storage.NewVisitedDbRepo(db)

	data := []*storage.DbCreateVisited{
		{
			Id:        1,
			Title:     "t1",
			PosterUrl: new("url1"),
			Status:    1,
		},
		{
			Id:        2,
			Title:     "t2",
			PosterUrl: new("url2"),
			Status:    2,
		},
	}

	// make expected list of DbVisited
	expected := make([]*storage.DbVisited, 0, len(data))
	for _, v := range data {
		expected = append(expected, &storage.DbVisited{
			Id:        v.Id,
			Title:     v.Title,
			PosterUrl: v.PosterUrl,
			Status:    v.Status,
		})
	}

	// run method
	err := repo.AddBatch(data)
	assert.Nil(t, err)

	// get list
	got, err := repo.List(-1, 0)

	// unset AddDate for comparsion
	for _, gotElem := range got.Items {
		gotElem.AddDate = ""
	}

	// expect no err
	assert.Nil(t, err)

	// expect elems match
	assert.ElementsMatch(t, expected, got.Items)

	// CASE add batch with existed id
	elem := &storage.DbCreateVisited{
		Id:        1,
		Title:     "new title",
		PosterUrl: new("new poster url"),
	}

	err = repo.AddBatch([]*storage.DbCreateVisited{elem})
	assert.Nil(t, err)

	// get elem from db
	got2, err := repo.Get(1)
	assert.Nil(t, err)

	got2.AddDate = ""
	assert.EqualValues(t, expected[0], got2)
}

func TestVisitedRepo_List(t *testing.T) {
	db := storage.TestInMemoDb(t)
	defer db.Close()

	repo := storage.NewVisitedDbRepo(db)

	DATA_LEN := 10
	data := make([]*storage.DbCreateVisited, 0, DATA_LEN)
	for i := 1; i <= 10; i++ {
		data = append(data, &storage.DbCreateVisited{
			Id:        i,
			Title:     "t" + strconv.Itoa(i),
			PosterUrl: new("url" + strconv.Itoa(i)),
			Status:    1,
		})
	}

	makeExpected := func(data []*storage.DbCreateVisited) []*storage.DbVisited {
		expected := make([]*storage.DbVisited, 0, len(data))
		for _, v := range data {
			expected = append(expected, &storage.DbVisited{
				Id:        v.Id,
				Title:     v.Title,
				PosterUrl: v.PosterUrl,
				Status:    v.Status,
			})
		}
		return expected
	}

	// add test data
	err := repo.AddBatch(data)
	assert.Nil(t, err)

	// test get all
	expected := makeExpected(data)

	got, err := repo.List(-1, 0)

	assert.Nil(t, err)
	assert.Equal(t, DATA_LEN, got.TotalCount)
	assert.Equal(t, DATA_LEN, got.Count)

	clearAddDate(got.Items)

	assert.ElementsMatch(t, expected, got.Items)

	// test limit offset
	// first part
	got, err = repo.List(6, 0)
	assert.Nil(t, err)
	assert.Equal(t, DATA_LEN, got.TotalCount)
	assert.Equal(t, 6, got.Count)
	clearAddDate(got.Items)
	assert.ElementsMatch(t, expected[:6], got.Items)

	// last part
	got, err = repo.List(6, 6)
	assert.Nil(t, err)
	assert.Equal(t, DATA_LEN, got.TotalCount)
	assert.Equal(t, 4, got.Count)
	clearAddDate(got.Items)
	assert.ElementsMatch(t, expected[6:], got.Items)
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
		PosterUrl: new("url1"),
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
