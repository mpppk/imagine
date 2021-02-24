package model

import (
	"fmt"
)

type TagID uint64
type UnindexedTag struct {
	ID   TagID  `json:"id"`
	Name string `json:"name"`
}

func (t *UnindexedTag) Index(index int) (*TagWithIndex, error) {
	return NewTagWithIndex(t.ID, t.Name, index)
}

// NewTag construct and returns UnindexedTag
func NewTag(id TagID, name string) (*UnindexedTag, error) {
	if name == "" {
		return nil, fmt.Errorf("failed to create tag: name is empty")
	}
	return &UnindexedTag{ID: id, Name: name}, nil
}

func (t *UnindexedTag) Unregister() *UnregisteredTag {
	return &UnregisteredTag{Name: t.Name}
}

func (t *UnindexedTag) SafeUnregister() (*UnregisteredTag, error) {
	if t.ID != 0 {
		return nil, fmt.Errorf("failed to unregister tag because it has ID(%d)", t.ID)
	}
	return t.Unregister(), nil
}

type UnregisteredTag struct {
	Name string
}

func (t *UnregisteredTag) Register(id TagID) *UnindexedTag {
	return &UnindexedTag{
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

func (t *UnindexedTag) GetID() uint64 {
	return uint64(t.ID)
}

func (t *UnindexedTag) SetID(id uint64) {
	t.ID = TagID(id)
}

type TagWithIndex struct {
	*UnindexedTag
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
	return &TagWithIndex{UnindexedTag: tag, Index: index}, nil
}

func (t *TagWithIndex) Unregister() *UnregisteredTagWithIndex {
	return &UnregisteredTagWithIndex{
		UnregisteredTag: t.UnindexedTag.Unregister(),
		Index:           t.Index,
	}
}

func (t *TagWithIndex) SafeUnregister() (*UnregisteredTagWithIndex, error) {
	if t.ID != 0 {
		return nil, fmt.Errorf("failed to unregister tag because it has ID(%d)", t.ID)
	}
	return t.Unregister(), nil
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
		UnindexedTag: t.UnregisteredTag.Register(id),
		Index:        t.Index,
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

// SubSetByID returns sub TagSet. new TagSet contains tag which have either provided ID.
func (t *TagSet) SubSetByID(idList []TagID) *TagSet {
	m := map[TagID]struct{}{}
	for _, id := range idList {
		m[id] = struct{}{}
	}
	return t.SubSetBy(func(tag *TagWithIndex) bool {
		_, ok := m[tag.ID]
		return ok
	})
}

// SubSetByNames returns sub TagSet. new TagSet contains tag which have either provided names.
func (t *TagSet) SubSetByNames(names []string) *TagSet {
	m := map[string]struct{}{}
	for _, name := range names {
		m[name] = struct{}{}
	}
	return t.SubSetBy(func(tag *TagWithIndex) bool {
		_, ok := m[tag.Name]
		return ok
	})
}

// SplitByNames splits TagSet to two sub sets based on provided tag names.
func (t *TagSet) SplitByNames(names []string) (existsTagSet, nonExistsTagSet *TagSet) {
	m := map[string]struct{}{}
	for _, name := range names {
		m[name] = struct{}{}
	}
	return t.SplitBy(func(tag *TagWithIndex) bool {
		_, ok := m[tag.Name]
		return ok
	})
}

// SplitByID splits TagSet to two sub sets based on provided tag ID.
func (t *TagSet) SplitByID(idList []TagID) (existsTagSet, nonExistsTagSet *TagSet) {
	m := map[TagID]struct{}{}
	for _, id := range idList {
		m[id] = struct{}{}
	}
	return t.SplitBy(func(tag *TagWithIndex) bool {
		_, ok := m[tag.ID]
		return ok
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
