package storage

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type VisitedDbRepo struct {
	db *sqlx.DB
}

func NewVisitedDbRepo(db *sqlx.DB) *VisitedDbRepo {
	repo := &VisitedDbRepo{
		db: db,
	}
	return repo
}

func (r *VisitedDbRepo) List() ([]*DbVisited, error) {
	visiteds := []*DbVisited{}

	err := r.db.Select(&visiteds, "select * from visited")

	return visiteds, err
}

func (r *VisitedDbRepo) Add(data *DbCreateVisited) (*DbVisited, error) {
	result, err := r.db.Exec("insert into visited(id, title, poster_url, status) values(?,?,?,?)", data.Id, data.Title, data.PosterUrl, data.Status)
	if err != nil {
		return nil, err
	}

	rowId, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return r.Get(int(rowId))
}

func (r *VisitedDbRepo) AddBatch(data []*DbCreateVisited) error {
	_, err := r.db.NamedExec("insert or ignore into visited(id, title, poster_url, status) values(:id,:title,:poster_url,:status)", data)

	return err
}

func (r *VisitedDbRepo) Get(id int) (*DbVisited, error) {
	visited := &DbVisited{}

	err := r.db.Get(visited, "select * from visited where id=$1", id)

	return visited, err
}

func (r *VisitedDbRepo) UpdatePartial(id int, updateData map[DbVisitedFieldName]interface{}) error {
	updateFields := make([]string, 0, len(updateData))
	values := make([]interface{}, 0, len(updateData))
	values = append(values, id)

	idx := 2
	for key, val := range updateData {
		updateFields = append(updateFields, fmt.Sprintf("%s=$%d", key, idx))
		values = append(values, val)
		idx++
	}

	_, err := r.db.Exec("UPDATE visited SET "+strings.Join(updateFields, ",")+" WHERE id=$1", values...)
	return err
}

func (r *VisitedDbRepo) Delete(id int) error {
	result, err := r.db.Exec("DELETE FROM visited WHERE id=$1", id)
	if err != nil {
		return err
	}

	rowsCount, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsCount == 0 {
		return sql.ErrNoRows
	}

	return err
}
