package repoimpl

import (
	"encoding/json"
	"fmt"

	"github.com/mpppk/imagine/domain/model"
	bolt "go.etcd.io/bbolt"
)

type BBoltBaseTag struct {
	bucketName string
	base       *boltRepository
}

func NewBBoltBaseTag(b *bolt.DB, bucketName string) *BBoltBaseTag {
	return &BBoltBaseTag{
		bucketName: bucketName,
		base:       newBoltRepository(b),
	}
}

func (b *BBoltBaseTag) createBucketNames(ws model.WSName) []string {
	return []string{string(ws), b.bucketName}
}

func (b *BBoltBaseTag) loBucketFunc(ws model.WSName, f func(bucket *bolt.Bucket) error) error {
	return b.base.loBucketFunc(b.createBucketNames(ws), f)
}

func (b *BBoltBaseTag) count(ws model.WSName) (int, error) {
	tagSet, err := b.ListAsSet(ws)
	if err != nil {
		return 0, fmt.Errorf("failed to get tag set: %w", err)
	}

	tagMap, _ := tagSet.ToMap()
	return len(tagMap), nil
}

func (b *BBoltBaseTag) Add(ws model.WSName, unregisteredTag *model.UnregisteredTag) (*model.TagWithIndex, error) {
	errMsg := "failed to add tag"
	tagNum, err := b.count(ws)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errMsg, err)
	}

	unregisteredTagWithIndex, err := model.NewUnregisteredTagWithIndexFromUnregisteredTag(unregisteredTag, tagNum)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errMsg, err)
	}
	tag, err := b.AddWithIndex(ws, unregisteredTagWithIndex)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errMsg, err)
	}
	return tag, nil
}

func (b *BBoltBaseTag) AddWithIndex(ws model.WSName, unregisteredTagWithIndex *model.UnregisteredTagWithIndex) (*model.TagWithIndex, error) {
	errMsg := "failed to add tag with index"
	id, err := b.base.AddWithID(b.createBucketNames(ws), unregisteredTagWithIndex.Register(0))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errMsg, err)
	}
	return unregisteredTagWithIndex.Register(model.TagID(id)), nil
}

func (b *BBoltBaseTag) AddByName(ws model.WSName, tagName string) (*model.TagWithIndex, bool, error) {
	errMsg := "failed to add tag to db by name"
	count, err := b.count(ws)
	if err != nil {
		return nil, false, fmt.Errorf("%s: %w", errMsg, err)
	}
	unregisteredTag, err := model.NewUnregisteredTagWithIndex(tagName, count)
	if err != nil {
		return nil, false, fmt.Errorf("%s: %w", errMsg, err)
	}

	tag, err := b.AddWithIndex(ws, unregisteredTag)
	if err != nil {
		return nil, false, fmt.Errorf("%s: %w", errMsg, err)
	}
	return tag, true, err
}

// AddByNames adds tags which have provided names. Then returns ID list of created tags.
// If tag which have same name exists, do nothing and return the exist tag ID.
// For example, assume that ["tag1", "tag2", "tag3"] are provided as tagNames, and "tag2" already exist with ID=1.
// In this case, return values is [2,1,3].
func (b *BBoltBaseTag) AddByNames(ws model.WSName, tagNames []string) ([]*model.TagWithIndex, error) {
	errMsg := "failed to add tags by names"
	tagSet, err := b.ListAsSet(ws)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errMsg, err)
	}

	tagMap, tagNameMap := tagSet.ToMap()
	lastIndex := len(tagMap)
	var tags []*model.TagWithIndex
	for _, name := range tagNames {
		if tag, ok := tagNameMap[name]; ok {
			tags = append(tags, tag)
			continue
		}
		tag, err := model.NewUnregisteredTagWithIndex(name, lastIndex)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", errMsg, err)
		}
		lastIndex++
		id, err := b.AddWithIndex(ws, tag)
		if err != nil {
			return nil, fmt.Errorf("failed to add tag by names: %w", err)
		}
		tags = append(tags, id)
	}
	return tags, err
}

func (b *BBoltBaseTag) Get(ws model.WSName, id model.TagID) (tagWithIndex *model.TagWithIndex, exist bool, err error) {
	data, exist, err := b.base.Get(b.createBucketNames(ws), uint64(id))
	if err != nil {
		return nil, exist, err
	}
	if !exist {
		return nil, false, nil
	}

	var a model.TagWithIndex
	if err := json.Unmarshal(data, &a); err != nil {
		return nil, exist, fmt.Errorf("failed to unmarshal json to tag. contents: %s: %w", string(data), err)
	}
	return &a, exist, nil
}

func (b *BBoltBaseTag) RecreateBucket(ws model.WSName) error {
	return b.base.recreateBucket(b.createBucketNames(ws))
}

// Save saves tag to bolt.
// If a tag with the same ID is already exists, update it by provided tag.
// If tag does not exist yet, add provided tag.
// If tag which has same name exists on tag history but different ID, return error.
func (b *BBoltBaseTag) Save(ws model.WSName, tagWithIndex *model.TagWithIndex) (*model.TagWithIndex, error) {
	id, err := b.base.SaveByID(b.createBucketNames(ws), tagWithIndex)
	return tagWithIndex.ReRegister(model.TagID(id)), err
}

func (b *BBoltBaseTag) ListByAsync(ws model.WSName, f func(tagWithIndex *model.TagWithIndex) bool, cap int) (assetChan <-chan *model.TagWithIndex, err error) {
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

func (b *BBoltBaseTag) ListAll(ws model.WSName) (assets []*model.TagWithIndex, err error) {
	return b.ListBy(ws, func(tag *model.TagWithIndex) bool { return true })
}

func (b *BBoltBaseTag) ListBy(ws model.WSName, f func(tag *model.TagWithIndex) bool) (assets []*model.TagWithIndex, err error) {
	eachF := func(tagWithIndex *model.TagWithIndex) error {
		if f(tagWithIndex) {
			assets = append(assets, tagWithIndex)
		}
		return nil
	}
	if err := b.ForEach(ws, eachF); err != nil {
		return nil, fmt.Errorf("failed to list tags: %w", err)
	}
	return
}

func (b *BBoltBaseTag) ListAsSet(ws model.WSName) (set *model.TagSet, err error) {
	tags, err := b.ListAll(ws)
	if err != nil {
		return nil, err
	}

	set = model.NewTagSet(nil)
	for _, tag := range tags {
		set.Set(tag)
	}
	return
}

func (b *BBoltBaseTag) ForEach(ws model.WSName, f func(tagWithIndex *model.TagWithIndex) error) error {
	return b.loBucketFunc(ws, func(bucket *bolt.Bucket) error {
		return bucket.ForEach(func(k, v []byte) error {
			var tagWithIndex model.TagWithIndex
			if err := json.Unmarshal(v, &tagWithIndex); err != nil {
				return fmt.Errorf("failed to unmarshal tag: %w", err)
			}
			return f(&tagWithIndex)
		})
	})
}

// Delete deletes tags which have provided ID.
// Internally, even if tag is deleted, it still reserved on tag bucket with `deleted` flag.
func (b *BBoltBaseTag) Delete(ws model.WSName, idList []model.TagID) error {
	// tag bucketから消す
	// deleted tag bucketに追加
	return nil
}
