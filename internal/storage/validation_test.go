package storage

import (
	"testing"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/stretchr/testify/assert"
)

func TestDBCreateVisited_Validation(t *testing.T) {
	test := []struct {
		name     string
		in       *DbCreateVisited
		err      error
		testFunc func(t *testing.T, in *DbCreateVisited, err validation.Errors)
	}{
		{
			"ok",
			&DbCreateVisited{
				Id:        1,
				Title:     "title",
				PosterUrl: new("http://url"),
				Status:    1,
			},
			nil,
			nil,
		},
		{
			"id required",
			&DbCreateVisited{
				Id:        0,
				Title:     "title",
				PosterUrl: new("http://url"),
				Status:    1,
			},
			validation.Errors{},
			func(t *testing.T, in *DbCreateVisited, err validation.Errors) {
				assert.Contains(t, err, "id")
			},
		},
		{
			"title required",
			&DbCreateVisited{
				Id:        1,
				Title:     "",
				PosterUrl: new("http://url"),
				Status:    1,
			},
			validation.Errors{},
			func(t *testing.T, in *DbCreateVisited, err validation.Errors) {
				assert.Contains(t, err, "title")
			},
		},
		{
			"poster url empty",
			&DbCreateVisited{
				Id:     1,
				Title:  "title",
				Status: 1,
			},
			nil,
			nil,
		},
		{
			"poster url invalid",
			&DbCreateVisited{
				Id:        1,
				Title:     "",
				PosterUrl: new("some"),
				Status:    1,
			},
			validation.Errors{},
			func(t *testing.T, in *DbCreateVisited, err validation.Errors) {
				assert.Contains(t, err, "poster_url")
			},
		},
	}

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.in.Validate()

			if tt.err != nil {
				assert.ErrorAs(t, err, &tt.err)

				if tt.testFunc != nil {
					tt.testFunc(t, tt.in, tt.err.(validation.Errors))
				}
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestDBPatchVisited_Validation(t *testing.T) {
	tests := []struct {
		name     string
		in       *DbPatchVisited
		err      error
		testFunc func(t *testing.T, in *DbPatchVisited, err validation.Errors)
	}{
		{
			"ok",
			&DbPatchVisited{
				Title:     new("title"),
				PosterUrl: new("http://url"),
				Status:    new(1),
			},
			nil,
			nil,
		},
		{
			"title empty",
			&DbPatchVisited{
				PosterUrl: new("http://url"),
				Status:    new(1),
			},
			nil,
			nil,
		},
		{
			"poster url empty",
			&DbPatchVisited{
				Title:     new("title"),
				PosterUrl: new(""),
				Status:    new(1),
			},
			nil,
			nil,
		},
		{
			"poster url invalid",
			&DbPatchVisited{
				PosterUrl: new("some"),
			},
			validation.Errors{},
			func(t *testing.T, in *DbPatchVisited, err validation.Errors) {
				assert.Contains(t, err, "poster_url")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.in.Validate()

			if tt.err != nil {
				assert.ErrorAs(t, err, &tt.err)

				if tt.testFunc != nil {
					tt.testFunc(t, tt.in, tt.err.(validation.Errors))
				}
			} else {
				assert.Nil(t, err)
			}
		})
	}
}
