package storage

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

type DbVisitedFieldName string

const (
    DbVisitedFieldNameTitle DbVisitedFieldName = "title"
    DbVisitedFieldNamePosterURL DbVisitedFieldName ="poster_url"
    DbVisitedFieldNameStatus DbVisitedFieldName = "status"
)
