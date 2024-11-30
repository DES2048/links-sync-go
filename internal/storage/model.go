package storage

type DbVisited struct {
	Id        int    `db:"id"`
	Title     string `db:"title"`
	PosterUrl string `db:"poster_url"`
	Status    int    `db:"status"`
}

type DbCreateVisited struct {
	Title     string `db:"title"`
	PosterUrl string `db:"poster_url"`
	Status    int    `db:"status"`
}
