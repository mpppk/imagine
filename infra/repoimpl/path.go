package repoimpl

import (
	"fmt"

	"github.com/mpppk/imagine/domain/model"
	"go.etcd.io/bbolt"
	bolt "go.etcd.io/bbolt"
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
		return 0, false, fmt.Errorf("failed to get asset ID by path. path: %s: %w", path, err)
	} else if !exist {
		return 0, false, nil
	}
	return model.AssetID(id2), exist, nil
}

// FilterExistPath returns paths which does not exist yet
func (p *bboltPathRepository) FilterExistPath(ws model.WSName, paths []string) (notExistPaths []string, err error) {
	idList, err := p.ListByPath(ws, paths)
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

// ListByPath lists asset ID of provided paths.
// If path does not exist in db, 0 is used as ID.
func (p *bboltPathRepository) ListByPath(ws model.WSName, paths []string) (idList []model.AssetID, err error) {
	dataList, err := p.base.multipleGetByString(createPathBucketNames(ws), paths)
	if err != nil {
		return nil, err
	}
	for _, data := range dataList {
		if data == nil {
			idList = append(idList, 0) // FIXME
		} else {
			idList = append(idList, model.AssetID(btoi(data)))
		}
	}
	return
}

func (p *bboltPathRepository) Add(ws model.WSName, path string, assetID model.AssetID) error {
	return p.base.addIDWithStringKey(createPathBucketNames(ws), path, uint64(assetID))
}

func (p *bboltPathRepository) AddList(ws model.WSName, paths []string, assetIDList []model.AssetID) error {
	return p.base.addIDListWithStringKey(createPathBucketNames(ws), paths, model.AssetIDListToUint64List(assetIDList))
}

func (p *bboltPathRepository) DeleteAll(ws model.WSName) error {
	return p.base.recreateBucket(createPathBucketNames(ws))
}

func (p *bboltPathRepository) ListAll(ws model.WSName) ([]string, []model.AssetID, error) {
	return p.ListBy(ws, func(p string, i model.AssetID) bool { return true })
}

func (p *bboltPathRepository) ListBy(ws model.WSName, f func(path string, id model.AssetID) bool) (paths []string, idList []model.AssetID, err error) {
	eachF := func(path string, id model.AssetID) error {
		if f(path, id) {
			paths = append(paths, path)
			idList = append(idList, id)
		}
		return nil
	}
	if err := p.ForEach(ws, eachF); err != nil {
		return nil, nil, fmt.Errorf("failed to list paths: %w", err)
	}
	return
}

func (p *bboltPathRepository) ForEach(ws model.WSName, f func(path string, id model.AssetID) error) error {
	return p.base.loBucketFunc(createPathBucketNames(ws), func(bucket *bolt.Bucket) error {
		return bucket.ForEach(func(k, v []byte) error {
			return f(string(k), model.AssetID(btoi(v)))
		})
	})
}
