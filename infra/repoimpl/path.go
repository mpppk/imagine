package repoimpl

import "go.etcd.io/bbolt"

type pathRepository struct {
	base *boltRepository
}

func newPathRepository(b *bbolt.DB) *pathRepository {
	return &pathRepository{
		base: newBoltRepository(b),
	}
}

//func (p *pathRepository) SaveIfNotExist(workSpace, path string) bool {
//
//}
