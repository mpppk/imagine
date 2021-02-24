package repoimpl

import (
	"fmt"

	"github.com/mpppk/imagine/domain/repository"
	"github.com/mpppk/imagine/infra/blt"

	"github.com/mpppk/imagine/domain/model"
	bolt "go.etcd.io/bbolt"
)

type BBoltTag struct {
	*BBoltBaseTag
	historyRepository *BBoltBaseTag
	base              *blt.Repository
}

func NewBBoltTag(b *bolt.DB) repository.Tag {
	return &BBoltTag{
		BBoltBaseTag:      NewBBoltBaseTag(b, blt.TagBucketName),
		historyRepository: NewBBoltBaseTag(b, blt.TagHistoryBucketName),
		base:              blt.NewRepository(b),
	}
}

func (b *BBoltTag) Init(ws model.WSName) error {
	if err := b.base.CreateBucketIfNotExist(blt.CreateTagBucketNames(ws)); err != nil {
		return fmt.Errorf("failed to create tag bucket: %w", err)
	}
	if err := b.base.CreateBucketIfNotExist(blt.CreateTagHistoryBucketNames(ws)); err != nil {
		return fmt.Errorf("failed to create tag history bucket: %w", err)
	}
	return nil
}

func (b *BBoltTag) AddWithIndex(ws model.WSName, tagWithIndex *model.UnregisteredTag) (*model.Tag, error) {
	errMsg := "failed to add tag with index"

	//b.historyRepository.

	id, err := b.BBoltBaseTag.AddWithIndex(ws, tagWithIndex)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errMsg, err)
	}
	return id, nil
}

func (b *BBoltTag) AddByName(ws model.WSName, tagName string) (*model.Tag, bool, error) {
	errMsg := "failed to add tag by name"
	id, ok, err := b.BBoltBaseTag.AddByName(ws, tagName)
	if err != nil {
		return nil, false, fmt.Errorf("%s: %w", errMsg, err)
	}
	// FIXME: history
	return id, ok, nil
}

// AddByNames adds tags which have provided names. Then returns ID list of created tags.
// If tag which have same name exists, do nothing and return the exist tag ID.
// For example, assume that ["tag1", "tag2", "tag3"] are provided as tagNames, and "tag2" already exist with ID=1.
// In this case, return values is [2,1,3].
func (b *BBoltTag) AddByNames(ws model.WSName, tagNames []string) ([]*model.Tag, error) {
	errMsg := "failed to add tag by name"
	tags, err := b.BBoltBaseTag.AddByNames(ws, tagNames)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errMsg, err)
	}
	// FIXME: history
	return tags, nil
}
