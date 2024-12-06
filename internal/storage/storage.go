package storage

import (
	"links-sync-go/internal/config"

	"github.com/jmoiron/sqlx"
)

type IVisitedRepo interface {
	List() ([]*DbVisited, error)
	ListByIds(ids []int) ([]*DbVisited, error)
	Add(data *DbCreateVisited) (*DbVisited, error)
	AddBatch(data []*DbCreateVisited) error
	Get(id int) (*DbVisited, error)
	UpdatePartial(id int, updateData *DbPatchVisited) error
	Delete(id int) error
}

type IStorage interface {
	VisitedRepo() IVisitedRepo
}

type DbStorage struct {
	config      *config.Config
	db          *sqlx.DB
	visitedRepo *VisitedDbRepo
}

func NewDbStorage(config *config.Config) (*DbStorage, error) {
	db, err := sqlx.Open("sqlite3", config.Db.Url)
	if err != nil {
		return nil, err
	}

	schema := `
		CREATE TABLE IF NOT EXISTS visited (
			id integer PRIMARY KEY,
			title VARCHAR NOT NULL,
            add_date VARCHAR NOT NULL DEFAULT(datetime()),
			poster_url varchar,
			status integer NOT NULL DEFAULT 1
		);
	`

	db.MustExec(schema)

	visitedRepo := NewVisitedDbRepo(db)

	return &DbStorage{
		config:      config,
		db:          db,
		visitedRepo: visitedRepo,
	}, nil
}

func (s *DbStorage) VisitedRepo() IVisitedRepo {
	return s.visitedRepo
}
