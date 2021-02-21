package model

import (
	"fmt"
)

type TagID uint64
type Tag struct {
	ID   TagID  `json:"id"`
	Name string `json:"name"`
}

// NewTag construct and returns Tag
func NewTag(id TagID, name string) (*Tag, error) {
	if name == "" {
		return nil, fmt.Errorf("failed to create tag: name is empty")
	}
	return &Tag{ID: id, Name: name}, nil
}

func (t *Tag) Unregister() *UnregisteredTag {
	return &UnregisteredTag{Name: t.Name}
}

type UnregisteredTag struct {
	Name string
}

func (t *UnregisteredTag) Register(id TagID) *Tag {
	return &Tag{
		ID:   id,
		Name: t.Name,
	}
}

// NewUnregisteredTag construct and returns UnregisteredTag
func NewUnregisteredTag(name string) (*UnregisteredTag, error) {
	if name == "" {
		return nil, fmt.Errorf("failed to create tag: name is empty")
	}
	return &UnregisteredTag{Name: name}, nil
}

func (t *Tag) GetID() uint64 {
	return uint64(t.ID)
}

func (t *Tag) SetID(id uint64) {
	t.ID = TagID(id)
}

type TagWithIndex struct {
	*Tag
	Index int
}

// NewTagWithIndex construct and return TagWithIndex
func NewTagWithIndex(id TagID, name string, index int) (*TagWithIndex, error) {
	errMsg := "failed to create TagWithIndex"
	if index < 0 {
		return nil, fmt.Errorf("%s: negative index is not allowed(%d)", errMsg, index)
	}
	tag, err := NewTag(id, name)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errMsg, err)
	}
	return &TagWithIndex{Tag: tag, Index: index}, nil
}

func (t *TagWithIndex) Unregister() *UnregisteredTagWithIndex {
	return &UnregisteredTagWithIndex{
		UnregisteredTag: t.Tag.Unregister(),
		Index:           t.Index,
	}
}

func (t *TagWithIndex) ReRegister(id TagID) *TagWithIndex {
	newTag := *t
	newTag.ID = id
	return &newTag
}

type UnregisteredTagWithIndex struct {
	*UnregisteredTag
	Index int
}

func (t *UnregisteredTagWithIndex) Register(id TagID) *TagWithIndex {
	return &TagWithIndex{
		Tag:   t.UnregisteredTag.Register(id),
		Index: t.Index,
	}
}

// NewUnregisteredTagWithIndex construct and return UnregisteredTagWithIndex
func NewUnregisteredTagWithIndex(name string, index int) (*UnregisteredTagWithIndex, error) {
	errMsg := "failed to create TagWithIndex"
	if index < 0 {
		return nil, fmt.Errorf("%s: negative index is not allowed(%d)", errMsg, index)
	}
	tag, err := NewUnregisteredTag(name)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errMsg, err)
	}
	return &UnregisteredTagWithIndex{UnregisteredTag: tag, Index: index}, nil
}

// NewUnregisteredTagWithIndexFromUnregisteredTag construct and return UnregisteredTagWithIndex from UnregisteredTag
func NewUnregisteredTagWithIndexFromUnregisteredTag(tag *UnregisteredTag, index int) (*UnregisteredTagWithIndex, error) {
	errMsg := "failed to create TagWithIndex"
	if index < 0 {
		return nil, fmt.Errorf("%s: negative index is not allowed(%d)", errMsg, index)
	}
	return &UnregisteredTagWithIndex{UnregisteredTag: tag, Index: index}, nil
}

type TagSet struct {
	m     map[TagID]*TagWithIndex
	nameM map[string]*TagWithIndex
}

// NewTagSet returns TagSet.
// if nil is provided as tags, empty TagSet will be return.
func NewTagSet(tags []*TagWithIndex) *TagSet {
	tagSet := &TagSet{
		m:     map[TagID]*TagWithIndex{},
		nameM: map[string]*TagWithIndex{},
	}
	for _, tag := range tags {
		tagSet.Set(tag)
	}
	return tagSet
}

// Set set tag. If provided tag does not exists yet, set the tag and return true.
// If provided tag name already exists on TagSet and has different ID, do nothing and return false.
// Otherwise, add or update the tag, and return true.
// If provided tag already exists but has different ID, update the tag and return true.
// Otherwise do nothing and returns false.
func (t *TagSet) Set(tag *TagWithIndex) bool {
	if sameNamedTag, ok := t.nameM[tag.Name]; ok && sameNamedTag.ID != tag.ID {
		return false
	}
	t.m[tag.ID] = tag
	t.nameM[tag.Name] = tag
	return true
}

func (t *TagSet) Get(id TagID) (*TagWithIndex, bool) {
	tag, ok := t.m[id]
	return tag, ok
}

func (t *TagSet) GetByName(name string) (*TagWithIndex, bool) {
	tag, ok := t.nameM[name]
	return tag, ok
}

func (t *TagSet) SubSetBy(f func(tag *TagWithIndex) bool) *TagSet {
	subset := NewTagSet(nil)
	for _, tag := range t.m {
		if f(tag) {
			subset.Set(tag)
		}
	}
	return subset
}

// SplitBy splits TagSet to two sub sets based on return value of provided function.
func (t *TagSet) SplitBy(f func(tag *TagWithIndex) bool) (trueTagSet, falseTagSet *TagSet) {
	trueSet := NewTagSet(nil)
	falseSet := NewTagSet(nil)
	for _, tag := range t.m {
		if f(tag) {
			trueSet.Set(tag)
		} else {
			falseSet.Set(tag)
		}
	}
	return trueSet, falseSet
}

// SubSetByNames returns sub TagSet. new TagSet contains tag which have either provided names.
func (t *TagSet) SubSetByNames(names []string) *TagSet {
	return t.SubSetBy(func(tag *TagWithIndex) bool {
		for _, name := range names {
			if tag.Name == name {
				return true
			}
		}
		return false
	})
}

// SplitByNames splits TagSet to two sub sets based on provided tag names.
func (t *TagSet) SplitByNames(names []string) (existsTagSet, nonExistsTagSet *TagSet) {
	return t.SplitBy(func(tag *TagWithIndex) bool {
		for _, name := range names {
			if tag.Name == name {
				return true
			}
		}
		return false
	})
}

// ToMap returns two maps.
func (t *TagSet) ToMap() (map[TagID]*TagWithIndex, map[string]*TagWithIndex) {
	return t.m, t.nameM
}

// ToTags returns tags
func (t *TagSet) ToTags() (tags []*TagWithIndex) {
	for _, tag := range t.m {
		tags = append(tags, tag)
	}
	return
}
