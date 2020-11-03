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

func (p *bboltPathRepository) AddIfNotExist(ws model.WSName, path string, assetID model.AssetID) (added bool, err error) {
	_, exist, err := p.base.getByString(createPathBucketNames(ws), path)
	if err != nil {
		return false, err
	}

	if exist {
		return false, nil
	}

	if err := p.base.addWithStringKey(createPathBucketNames(ws), path, assetID); err != nil {
		return false, err
	}

	return true, nil
}
