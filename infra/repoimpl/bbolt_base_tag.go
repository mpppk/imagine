package repoimpl

import (
	"encoding/json"
	"fmt"

	"github.com/mpppk/imagine/domain/model"
	"github.com/mpppk/imagine/infra/blt"
	bolt "go.etcd.io/bbolt"
)

type BBoltBaseTag struct {
	bucketName string
	base       *blt.Repository
}

func NewBBoltBaseTag(b *bolt.DB, bucketName string) *BBoltBaseTag {
	return &BBoltBaseTag{
		bucketName: bucketName,
		base:       blt.NewRepository(b),
	}
}

func (b *BBoltBaseTag) createBucketNames(ws model.WSName) []string {
	return []string{string(ws), b.bucketName}
}

func (b *BBoltBaseTag) loBucketFunc(ws model.WSName, f func(bucket *bolt.Bucket) error) error {
	return b.base.LoBucketFunc(b.createBucketNames(ws), f)
}

func (b *BBoltBaseTag) bucketFunc(ws model.WSName, f func(bucket *bolt.Bucket) error) error {
	return b.base.BucketFunc(b.createBucketNames(ws), f)
}

func (b *BBoltBaseTag) count(ws model.WSName) (int, error) {
	tagSet, err := b.ListAsSet(ws)
	if err != nil {
		return 0, fmt.Errorf("failed to get tag set: %w", err)
	}

	tagMap, _ := tagSet.ToMap()
	return len(tagMap), nil
}

func (b *BBoltBaseTag) Add(ws model.WSName, unregisteredTag *model.UnregisteredUnindexedTag) (*model.Tag, error) {
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

func (b *BBoltBaseTag) AddWithIndex(ws model.WSName, unregisteredTagWithIndex *model.UnregisteredTag) (*model.Tag, error) {
	errMsg := "failed to add tag with index"
	id, err := b.base.AddWithID(b.createBucketNames(ws), unregisteredTagWithIndex.Register(0))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errMsg, err)
	}
	return unregisteredTagWithIndex.Register(model.TagID(id)), nil
}

func (b *BBoltBaseTag) AddByName(ws model.WSName, tagName string) (*model.Tag, bool, error) {
	errMsg := "failed to add tag to db by name"
	count, err := b.count(ws)
	if err != nil {
		return nil, false, fmt.Errorf("%s: %w", errMsg, err)
	}
	unregisteredTag, err := model.NewUnregisteredTag(tagName, count)
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
func (b *BBoltBaseTag) AddByNames(ws model.WSName, tagNames []string) ([]*model.Tag, error) {
	errMsg := "failed to add tags by names"
	tagSet, err := b.ListAsSet(ws)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errMsg, err)
	}

	tagMap, tagNameMap := tagSet.ToMap()
	lastIndex := len(tagMap)
	var tags []*model.Tag
	for _, name := range tagNames {
		if tag, ok := tagNameMap[name]; ok {
			tags = append(tags, tag)
			continue
		}
		tag, err := model.NewUnregisteredTag(name, lastIndex)
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

func (b *BBoltBaseTag) Get(ws model.WSName, id model.TagID) (tagWithIndex *model.Tag, exist bool, err error) {
	data, exist, err := b.base.Get(b.createBucketNames(ws), uint64(id))
	if err != nil {
		return nil, exist, err
	}
	if !exist {
		return nil, false, nil
	}

	var a model.Tag
	if err := json.Unmarshal(data, &a); err != nil {
		return nil, exist, fmt.Errorf("failed to unmarshal json to tag. contents: %s: %w", string(data), err)
	}
	return &a, exist, nil
}

// Save saves tag to bolt.
// If a tag with the same ID is already exists, update it by provided tag.
// If tag does not exist yet, add provided tag.
// If tag which has same name exists on tag history but different ID, return error.
func (b *BBoltBaseTag) Save(ws model.WSName, tagWithIndex *model.Tag) (*model.Tag, error) {
	id, err := b.base.SaveByID(b.createBucketNames(ws), tagWithIndex)
	return tagWithIndex.ReRegister(model.TagID(id)), err
}

func (b *BBoltBaseTag) ListByAsync(ws model.WSName, f func(tagWithIndex *model.Tag) bool, cap int) (assetChan <-chan *model.Tag, err error) {
	c := make(chan *model.Tag, cap)
	ec := make(chan error, 1)
	f2 := f
	if f2 == nil {
		f2 = func(tagWithIndex *model.Tag) bool {
			return true
		}
	}
	eachF := func(tagWithIndex *model.Tag) error {
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

func (b *BBoltBaseTag) ListAll(ws model.WSName) (assets []*model.Tag, err error) {
	return b.ListBy(ws, func(tag *model.Tag) bool { return true })
}

func (b *BBoltBaseTag) ListBy(ws model.WSName, f func(tag *model.Tag) bool) (assets []*model.Tag, err error) {
	eachF := func(tagWithIndex *model.Tag) error {
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

func (b *BBoltBaseTag) ForEach(ws model.WSName, f func(tagWithIndex *model.Tag) error) error {
	return b.loBucketFunc(ws, func(bucket *bolt.Bucket) error {
		return bucket.ForEach(func(k, v []byte) error {
			var tagWithIndex model.Tag
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
	return b.bucketFunc(ws, func(bucket *bolt.Bucket) error {
		for _, id := range idList {
			if err := bucket.Delete(blt.Itob(uint64(id))); err != nil {
				return err
			}
		}
		return nil
	})
}
