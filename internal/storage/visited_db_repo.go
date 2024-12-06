package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
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

func (r *VisitedDbRepo) ListByIds(ids []int) ([]*DbVisited, error) {
	visiteds := []*DbVisited{}

	query, args, err := sqlx.In("SELECT id, status from visited WHERE id IN (?);", ids)
	if err != nil {
		return nil, err
	}

	err = r.db.Select(&visiteds, query, args...)

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
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}

		return nil, err
	}
	return visited, nil
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
		return ErrNotFound
	}

	return err
}

func (r *VisitedDbRepo) UpdatePartial(id int, updateData *DbPatchVisited) error {
	// fields for update
	updateFields := make([]string, 0)
	// corresponding values
	values := make([]interface{}, 0)

	// get fields and its vals from struct ignoring nil values
	v := reflect.ValueOf(updateData).Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		fld := v.Field(i)

		var val interface{}

		if fld.IsNil() {
			continue // ignore if nil
		} else {
			val = fld.Elem().Interface() // pointer deref
		}

		// get db field from struct
		colname := strings.Split(t.Field(i).Tag.Get("db"), ",")[0]

		updateFields = append(updateFields, fmt.Sprintf("%s=?", colname))
		values = append(values, val)
	}

	// apending id
	values = append(values, id)
	query := "UPDATE visited SET " + strings.Join(updateFields, ",") + " WHERE id=?"

	// fmt.Println(query)
	// fmt.Println(values...)
	result, err := r.db.Exec(query, values...)
	
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrNotFound
	}

	// fmt.Printf("rows affected %v\n", rowsAffected)
	return err
}
