package storage

import (
	"fmt"

	"links-sync-go/internal/config"

	"github.com/jmoiron/sqlx"
)

type IVisitedRepo interface {
	List(limit int, offset int) (*DBListMeta[*DbVisited], error)
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
	db, err := NewDB(config.Db.Url)
	if err != nil {
		return nil, fmt.Errorf("failed to open db with url: %s, cause: %w", config.Db.Url, err)
	}
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
