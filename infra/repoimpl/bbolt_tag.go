package repoimpl

import (
	"fmt"

	"github.com/mpppk/imagine/domain/repository"

	"github.com/mpppk/imagine/domain/model"
	bolt "go.etcd.io/bbolt"
)

type BBoltTag struct {
	*BBoltBaseTag
	historyRepository *BBoltBaseTag
	base              *boltRepository
}

func NewBBoltTag(b *bolt.DB) repository.Tag {
	return &BBoltTag{
		BBoltBaseTag:      NewBBoltBaseTag(b, tagBucketName),
		historyRepository: NewBBoltBaseTag(b, tagHistoryBucketName),
		base:              newBoltRepository(b),
	}
}

//func (b *BBoltTag) loBucketFunc(ws model.WSName, f func(bucket *bolt.Bucket) error) error {
//	return b.base.loBucketFunc(createTagBucketNames(ws), f)
//}

func (b *BBoltTag) Init(ws model.WSName) error {
	if err := b.base.createBucketIfNotExist(createTagBucketNames(ws)); err != nil {
		return fmt.Errorf("failed to create tag bucket: %w", err)
	}
	if err := b.base.createBucketIfNotExist(createTagHistoryBucketNames(ws)); err != nil {
		return fmt.Errorf("failed to create tag history bucket: %w", err)
	}
	return nil
}

func (b *BBoltTag) Add(ws model.WSName, tagWithIndex *model.TagWithIndex) (model.TagID, error) {
	errMsg := "failed to add tag"
	id, err := b.BBoltBaseTag.Add(ws, tagWithIndex)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", errMsg, err)
	}
	id, err = b.historyRepository.Save(ws, tagWithIndex)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", errMsg, err)
	}
	return id, nil
	//tag, ok, err := b.base.get(createTagHistoryBucketNames(ws), tagWithIndex.ID)
	//if err != nil {
	//	return 0, fmt.Errorf("failed to get tag history: %w", err)
	//}
	//id, err := b.base.add(createTagBucketNames(ws), tagWithIndex)
	//if _, err := b.base.saveByID(createTagHistoryBucketNames(ws), tagWithIndex); err != nil {
	//	return 0, fmt.Errorf("failed to update tag history: %w", err)
	//}
	//return model.TagID(id), err
}

func (b *BBoltTag) AddByName(ws model.WSName, tagName string) (model.TagID, bool, error) {
	errMsg := "failed to add tag by name"
	id, ok, err := b.BBoltBaseTag.AddByName(ws, tagName)
	if err != nil {
		return 0, false, fmt.Errorf("%s: %w", errMsg, err)
	}
	// FIXME: history
	return id, ok, nil
	//tagSet, err := b.ListAsSet(ws)
	//if err != nil {
	//	return 0, false, fmt.Errorf("failed to get tag set: %w", err)
	//}
	//
	//tagMap, tagNameMap := tagSet.ToMap()
	//lastIndex := len(tagMap)
	//if _, ok := tagNameMap[tagName]; ok {
	//	return 0, false, nil
	//}
	//lastIndex++
	//tag := &model.TagWithIndex{Tag: &model.Tag{Name: tagName}, Index: lastIndex}
	//id, err := b.Add(ws, tag)
	//if err != nil {
	//	return 0, false, fmt.Errorf("failed to add tag to db by name: %w", err)
	//}
	//return id, true, err
}

// AddByNames adds tags which have provided names. Then returns ID list of created tags.
// If tag which have same name exists, do nothing and return the exist tag ID.
// For example, assume that ["tag1", "tag2", "tag3"] are provided as tagNames, and "tag2" already exist with ID=1.
// In this case, return values is [2,1,3].
func (b *BBoltTag) AddByNames(ws model.WSName, tagNames []string) ([]model.TagID, error) {
	errMsg := "failed to add tag by name"
	idList, err := b.BBoltBaseTag.AddByNames(ws, tagNames)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errMsg, err)
	}
	// FIXME: history
	return idList, nil
	//tagSet, err := b.ListAsSet(ws)
	//if err != nil {
	//	return nil, fmt.Errorf("failed to get tag set: %w", err)
	//}
	//
	//tagMap, tagNameMap := tagSet.ToMap()
	//lastIndex := len(tagMap)
	//var idList []model.TagID
	//for _, name := range tagNames {
	//	if tag, ok := tagNameMap[name]; ok {
	//		idList = append(idList, tag.ID)
	//		continue
	//	}
	//	lastIndex++
	//	tag := &model.TagWithIndex{Tag: &model.Tag{Name: name}, Index: lastIndex}
	//	id, err := b.Add(ws, tag)
	//	if err != nil {
	//		return nil, fmt.Errorf("failed to add tag by names: %w", err)
	//	}
	//	idList = append(idList, id)
	//}
	//return idList, err
}

//func (b *BBoltTag) Get(ws model.WSName, id model.TagID) (tagWithIndex *model.TagWithIndex, exist bool, err error) {
//	data, exist, err := b.base.get(createTagBucketNames(ws), uint64(id))
//	if err != nil {
//		return nil, exist, err
//	}
//	if !exist {
//		return nil, false, nil
//	}
//
//	var a model.TagWithIndex
//	if err := json.Unmarshal(data, &a); err != nil {
//		return nil, exist, fmt.Errorf("failed to unmarshal json to tagWithIndex. contents: %s: %w", string(data), err)
//	}
//	return &a, exist, nil
//}

//func (b *BBoltTag) GetFromHistory(ws model.WSName, id model.TagID) (tag *model.Tag, exist bool, err error) {
//	data, exist, err := b.base.get(createTagHistoryBucketNames(ws), uint64(id))
//	if err != nil {
//		return nil, exist, err
//	}
//	if !exist {
//		return nil, false, nil
//	}
//
//	var t model.Tag
//	if err := json.Unmarshal(data, &t); err != nil {
//		return nil, exist, fmt.Errorf("failed to unmarshal json to tag. contents: %s: %w", string(data), err)
//	}
//	return &t, exist, nil
//}

//func (b *BBoltTag) RecreateBucket(ws model.WSName) error {
//	return b.base.recreateBucket(createTagBucketNames(ws))
//}

// Save saves tag to bolt.
// If a tag with the same ID is already exists, update it by provided tag.
// If tag does not exist yet, add provided tag.
// If tag which has same name exists on tag history but different ID, return error.
//func (b *BBoltTag) Save(ws model.WSName, tagWithIndex *model.TagWithIndex) (model.TagID, error) {
//	tag, ok, err := b.GetFromHistory(ws, tagWithIndex.ID)
//	if tag, ok, err := b.base.get(createTagHistoryBucketNames(ws), uint64(tagWithIndex.ID)); err != nil {
//		return 0, fmt.Errorf("failed to save tag history: %w", err)
//	} else if ok && tag != tagWithIndex.ID {
//		return 0, fmt.Errorf("invalid Tag ID. provided: %d, history: %d", tagWithIndex.ID, id)
//	}
//	id, err := b.base.saveByID(createTagBucketNames(ws), tagWithIndex)
//	return model.TagID(id), err
//}

//func (b *BBoltTag) ListByAsync(ws model.WSName, f func(tagWithIndex *model.TagWithIndex) bool, cap int) (assetChan <-chan *model.TagWithIndex, err error) {
//	c := make(chan *model.TagWithIndex, cap)
//	ec := make(chan error, 1)
//	f2 := f
//	if f2 == nil {
//		f2 = func(tagWithIndex *model.TagWithIndex) bool {
//			return true
//		}
//	}
//	eachF := func(tagWithIndex *model.TagWithIndex) error {
//		if f2(tagWithIndex) {
//			c <- tagWithIndex
//		}
//		return nil
//	}
//
//	go func() {
//		if err := b.ForEach(ws, eachF); err != nil {
//			ec <- fmt.Errorf("failed to list assets: %w", err)
//		}
//		close(c)
//		close(ec)
//	}()
//	return c, nil
//}

//func (b *BBoltTag) ListAll(ws model.WSName) (assets []*model.TagWithIndex, err error) {
//	return b.ListBy(ws, func(tag *model.TagWithIndex) bool { return true })
//}

//func (b *BBoltTag) ListBy(ws model.WSName, f func(tag *model.TagWithIndex) bool) (assets []*model.TagWithIndex, err error) {
//	eachF := func(tagWithIndex *model.TagWithIndex) error {
//		if f(tagWithIndex) {
//			assets = append(assets, tagWithIndex)
//		}
//		return nil
//	}
//	if err := b.ForEach(ws, eachF); err != nil {
//		return nil, fmt.Errorf("failed to list tags: %w", err)
//	}
//	return
//}

//func (b *BBoltTag) ListAsSet(ws model.WSName) (set *model.TagSet, err error) {
//	tags, err := b.ListAll(ws)
//	if err != nil {
//		return nil, err
//	}
//
//	set = model.NewTagSet(nil)
//	for _, tag := range tags {
//		set.Set(tag.Tag)
//	}
//	return
//}

//func (b *BBoltTag) ForEach(ws model.WSName, f func(tagWithIndex *model.TagWithIndex) error) error {
//	return b.loBucketFunc(ws, func(bucket *bolt.Bucket) error {
//		return bucket.ForEach(func(k, v []byte) error {
//			var tagWithIndex model.TagWithIndex
//			if err := json.Unmarshal(v, &tagWithIndex); err != nil {
//				return fmt.Errorf("failed to unmarshal tagWithIndex: %w", err)
//			}
//			return f(&tagWithIndex)
//		})
//	})
//}

// Delete deletes tags which have provided ID.
// Internally, even if tag is deleted, it still reserved on tag bucket with `deleted` flag.
//func (b *BBoltTag) Delete(ws model.WSName, idList []model.TagID) error {
//	// tag bucketから消す
//	// deleted tag bucketに追加
//	return nil
//}
