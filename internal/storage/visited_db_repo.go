package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/mattn/go-sqlite3"
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

	query, args, err := sqlx.In("SELECT * from visited WHERE id IN (?);", ids)
	if err != nil {
		return nil, err
	}

	err = r.db.Select(&visiteds, query, args...)

	return visiteds, err
}

func (r *VisitedDbRepo) Add(data *DbCreateVisited) (*DbVisited, error) {
	// set default status
	status := 1
	if data.Status > 0 {
		status = data.Status
	}
	result, err := r.db.Exec("insert into visited(id, title, poster_url, status) values(?,?,?,?)", data.Id, data.Title, data.PosterUrl, status)
	if err != nil {

		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) {
			if sqliteErr.Code == sqlite3.ErrConstraint && sqliteErr.ExtendedCode == sqlite3.ErrConstraintPrimaryKey {
				err = NewErrAlreadyExists(fmt.Sprintf("topic with id %d already exists", data.Id))
			}
		}

		return nil, err
	}

	rowID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return r.Get(int(rowID))
}

func (r *VisitedDbRepo) AddBatch(data []*DbCreateVisited) error {
	for _, d := range data {
		if d.Status == 0 {
			d.Status = 1
		}
	}
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
		// FIXME: potentional sql injection
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
