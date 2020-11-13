package repoimpl

import (
	"fmt"

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

func (p *bboltPathRepository) Get(ws model.WSName, path string) (id model.AssetID, exist bool, err error) {
	id2, exist, err := p.base.getIDByString(createPathBucketNames(ws), path)
	if err != nil {
		return 0, false, fmt.Errorf("failed to get asset id by path. path: %s: %w", path, err)
	} else if !exist {
		return 0, false, nil
	}
	return model.AssetID(id2), exist, nil
}

// FilterExistPath returns paths which does not exist yet
func (p *bboltPathRepository) FilterExistPath(ws model.WSName, paths []string) (notExistPaths []string, err error) {
	idList, err := p.GetList(ws, paths)
	if err != nil {
		return nil, err
	}

	for i, id := range idList {
		if id == 0 {
			notExistPaths = append(notExistPaths, paths[i])
		}
	}
	return
}

func (p *bboltPathRepository) GetList(ws model.WSName, paths []string) (idList []uint64, err error) {
	dataList, err := p.base.multipleGetByString(createPathBucketNames(ws), paths)
	if err != nil {
		return nil, err
	}
	for _, data := range dataList {
		if data == nil {
			idList = append(idList, 0) // FIXME
		} else {
			idList = append(idList, btoi(data))
		}
	}
	return
}

func (p *bboltPathRepository) Add(ws model.WSName, path string, assetID model.AssetID) error {
	return p.base.addIDWithStringKey(createPathBucketNames(ws), path, uint64(assetID))
}

func (p *bboltPathRepository) AddList(ws model.WSName, paths []string, assetIDList []model.AssetID) error {
	var dataList []interface{}
	for _, id := range assetIDList {
		dataList = append(dataList, id)
	}
	return p.base.addListWithStringKey(createPathBucketNames(ws), paths, dataList)
}

func (p *bboltPathRepository) DeleteAll(ws model.WSName) error {
	return p.base.recreateBucket(createPathBucketNames(ws))
}
