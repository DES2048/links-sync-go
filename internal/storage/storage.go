package storage

type IVisitedRepo interface {
	List() ([]*DbVisited, error)
	Add(data *DbCreateVisited) (DbVisited, error)
	Get(id int) (*DbVisited, error)
	UpdatePartial(id int, updateData map[string]interface{}) error
	Delete(id int) error
}

type IStorage interface {
	VisitedRepo() *IVisitedRepo
}
