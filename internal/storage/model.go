package storage

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

var (
	VisitedPosterUrlRule = []validation.Rule{is.URL}
	VisitedStatusRule    = []validation.Rule{validation.In(1, 2)}
	VisitedTitleRule = []validation.Rule{validation.Required}
)

type DbVisited struct {
	Id        int    `db:"id"`
	Title     string `db:"title"`
	AddDate   string `db:"add_date"`
	PosterUrl string `db:"poster_url"`
	Status    int    `db:"status"`
}

type DbCreateVisited struct {
	Id        int    `db:"id"`
	Title     string `db:"title"`
	PosterUrl string `db:"poster_url"`
	Status    int    `db:"status"`
}

func (o *DbCreateVisited) Validate() error {
	return validation.ValidateStruct(o,
		validation.Field(&o.Id, validation.Required),
		validation.Field(&o.Title, VisitedTitleRule...),
		validation.Field(&o.PosterUrl, VisitedPosterUrlRule...),
		validation.Field(&o.Status, VisitedStatusRule...),
	)
}

type DbPatchVisited struct {
	Title     *string `db:"title"`
	PosterUrl *string `db:"poster_url"`
	Status    *int    `db:"status"`
}

func (o *DbPatchVisited) Validate() error {
	return validation.ValidateStruct(o,
		validation.Field(&o.Title, VisitedTitleRule...),
		validation.Field(&o.PosterUrl, VisitedPosterUrlRule...),
		validation.Field(&o.Status, VisitedStatusRule...),
	)
}
