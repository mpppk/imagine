package repoimpl

import (
	"github.com/mpppk/imagine/domain/model"
	"go.etcd.io/bbolt"
)

type bboltPathRepository struct {
	base *boltRepository
}

func newBBoltPathRepository(b *bbolt.DB) *bboltPathRepository {
	return &bboltPathRepository{
		base: newBoltRepository(b),
	}
}

func (p *bboltPathRepository) Get(ws model.WSName, path string) (id uint64, exist bool, err error) {
	data, exist, err := p.base.getByString(createPathBucketNames(ws), path)
	if err != nil {
		return 0, false, err
	} else if !exist {
		return 0, false, nil
	}
	return btoi(data), exist, nil
}

func (p *bboltPathRepository) Add(ws model.WSName, path string, assetID model.AssetID) error {
	return p.base.addWithStringKey(createPathBucketNames(ws), path, assetID)
}
