package repoimpl

import (
	"encoding/json"
	"fmt"

	"github.com/mpppk/imagine/domain/repository"

	"github.com/mpppk/imagine/domain/model"
	bolt "go.etcd.io/bbolt"
)

type BBoltTag struct {
	base *boltRepository
}

func NewBBoltTag(b *bolt.DB) repository.Tag {
	return &BBoltTag{
		base: newBoltRepository(b),
	}
}

func (b *BBoltTag) loBucketFunc(ws model.WSName, f func(bucket *bolt.Bucket) error) error {
	return b.base.loBucketFunc(createTagBucketNames(ws), f)
}

func (b *BBoltTag) Init(ws model.WSName) error {
	return b.base.createBucketIfNotExist(createTagBucketNames(ws))
}

func (b *BBoltTag) Add(ws model.WSName, tagWithIndex *model.TagWithIndex) (model.TagID, error) {
	id, err := b.base.addByID(createTagBucketNames(ws), tagWithIndex)
	return model.TagID(id), err
}

func (b *BBoltTag) AddByName(ws model.WSName, tagName string) (model.TagID, bool, error) {
	tagSet, err := b.ListAsSet(ws)
	if err != nil {
		return 0, false, fmt.Errorf("failed to get tag set: %w", err)
	}

	tagMap, tagNameMap := tagSet.ToMap()
	lastIndex := len(tagMap)
	if _, ok := tagNameMap[tagName]; ok {
		return 0, false, nil
	}
	lastIndex++
	tag := model.TagWithIndex{Tag: &model.Tag{Name: tagName}, Index: lastIndex}
	id, err := b.base.addByID(createTagBucketNames(ws), tag)
	if err != nil {
		return 0, false, fmt.Errorf("failed to add tag to db: %w", err)
	}
	return model.TagID(id), true, err
}

func (b *BBoltTag) AddByNames(ws model.WSName, tagNames []string) ([]model.TagID, error) {
	tagSet, err := b.ListAsSet(ws)
	if err != nil {
		return nil, fmt.Errorf("failed to get tag set: %w", err)
	}

	tagMap, tagNameMap := tagSet.ToMap()
	lastIndex := len(tagMap)
	var idList []model.TagID
	for _, name := range tagNames {
		if _, ok := tagNameMap[name]; ok {
			continue
		}
		lastIndex++
		tag := model.TagWithIndex{Tag: &model.Tag{Name: name}, Index: lastIndex}
		id, err := b.base.addByID(createTagBucketNames(ws), tag)
		if err != nil {
			return nil, fmt.Errorf("failed to add tag to db: %w", err)
		}
		idList = append(idList, model.TagID(id))
	}
	return idList, err
}

func (b *BBoltTag) Get(ws model.WSName, id model.TagID) (tagWithIndex *model.TagWithIndex, exist bool, err error) {
	data, exist, err := b.base.get(createTagBucketNames(ws), uint64(id))
	if err != nil {
		return nil, exist, err
	}
	if !exist {
		return nil, false, nil
	}

	var a model.TagWithIndex
	if err := json.Unmarshal(data, &a); err != nil {
		return nil, exist, fmt.Errorf("failed to unmarshal json to tagWithIndex. contents: %s: %w", string(data), err)
	}
	return &a, exist, nil
}

func (b *BBoltTag) RecreateBucket(ws model.WSName) error {
	return b.base.recreateBucket(createTagBucketNames(ws))
}

func (b *BBoltTag) Update(ws model.WSName, tagWithIndex *model.TagWithIndex) error {
	return b.base.updateByID(createTagBucketNames(ws), tagWithIndex)
}

func (b *BBoltTag) ListByAsync(ws model.WSName, f func(tagWithIndex *model.TagWithIndex) bool, cap int) (assetChan <-chan *model.TagWithIndex, err error) {
	c := make(chan *model.TagWithIndex, cap)
	ec := make(chan error, 1)
	f2 := f
	if f2 == nil {
		f2 = func(tagWithIndex *model.TagWithIndex) bool {
			return true
		}
	}
	eachF := func(tagWithIndex *model.TagWithIndex) error {
		if f2(tagWithIndex) {
			c <- tagWithIndex
		}
		return nil
	}

	go func() {
		if err := b.ForEach(ws, eachF); err != nil {
			ec <- fmt.Errorf("failed to list assets: %w", err)
		}
		close(c)
		close(ec)
	}()
	return c, nil
}

func (b *BBoltTag) ListAll(ws model.WSName) (assets []*model.TagWithIndex, err error) {
	return b.ListBy(ws, func(tag *model.TagWithIndex) bool { return true })
}

func (b *BBoltTag) ListBy(ws model.WSName, f func(tag *model.TagWithIndex) bool) (assets []*model.TagWithIndex, err error) {
	eachF := func(tagWithIndex *model.TagWithIndex) error {
		if f(tagWithIndex) {
			assets = append(assets, tagWithIndex)
		}
		return nil
	}
	if err := b.ForEach(ws, eachF); err != nil {
		return nil, fmt.Errorf("failed to list assets: %w", err)
	}
	return
}

func (b *BBoltTag) ListAsSet(ws model.WSName) (set *model.TagSet, err error) {
	tags, err := b.ListAll(ws)
	if err != nil {
		return nil, err
	}

	set = model.NewTagSet(nil)
	for _, tag := range tags {
		set.Set(tag.Tag)
	}
	return
}

func (b *BBoltTag) ForEach(ws model.WSName, f func(tagWithIndex *model.TagWithIndex) error) error {
	return b.loBucketFunc(ws, func(bucket *bolt.Bucket) error {
		return bucket.ForEach(func(k, v []byte) error {
			var tagWithIndex model.TagWithIndex
			if err := json.Unmarshal(v, &tagWithIndex); err != nil {
				return fmt.Errorf("failed to unmarshal tagWithIndex: %w", err)
			}
			return f(&tagWithIndex)
		})
	})
}
