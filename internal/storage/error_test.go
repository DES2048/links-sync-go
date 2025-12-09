package storage_test

import (
	"errors"
	"testing"

	"links-sync-go/internal/storage"

	"github.com/stretchr/testify/assert"
)

func TestIsErrNotFound(t *testing.T) {
	// test ErrNotFound var
	err1 := storage.ErrNotFound

	assert.True(t, storage.IsErrNotFound(err1))

	// test err2
	err2 := storage.ErrNotFound2("visited not found")
	assert.True(t, storage.IsErrNotFound(err2))

	// test false
	err3 := errors.New("some error")
	assert.False(t, storage.IsErrNotFound(err3))
}
