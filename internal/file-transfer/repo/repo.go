package repo

import (
	"sync"

	"go.mongodb.org/mongo-driver/mongo"
)

var (
	once sync.Once
	S    *dataSource
)

type IRepo interface {
	NewTagRepo() MessageRepo
}

type dataSource struct {
	db *mongo.Client
}

var _ IRepo = (*dataSource)(nil)

func NewDataSource(db *mongo.Client) *dataSource {
	once.Do(func() {
		S = &dataSource{db}
	})
	return S
}

func (s *dataSource) NewTagRepo() MessageRepo {
	return newMessageRepo(s.db)
}
