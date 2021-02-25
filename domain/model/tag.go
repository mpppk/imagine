package model

import (
	"fmt"
)

type TagID uint64
type UnindexedTag struct {
	ID   TagID  `json:"id"`
	Name string `json:"name"`
}

func (t *UnindexedTag) Index(index int) (*Tag, error) {
	return NewTag(t.ID, t.Name, index)
}

// NewUnindexedTag construct and returns UnindexedTag
func NewUnindexedTag(id TagID, name string) (*UnindexedTag, error) {
	if name == "" {
		return nil, fmt.Errorf("failed to create tag: name is empty")
	}
	return &UnindexedTag{ID: id, Name: name}, nil
}

func (t *UnindexedTag) Unregister() *UnregisteredUnindexedTag {
	return &UnregisteredUnindexedTag{Name: t.Name}
}

func (t *UnindexedTag) SafeUnregister() (*UnregisteredUnindexedTag, error) {
	if t.ID != 0 {
		return nil, fmt.Errorf("failed to unregister tag because it has ID(%d)", t.ID)
	}
	return t.Unregister(), nil
}

type UnregisteredUnindexedTag struct {
	Name string
}

func (t *UnregisteredUnindexedTag) Register(id TagID) *UnindexedTag {
	return &UnindexedTag{
		ID:   id,
		Name: t.Name,
	}
}

// NewUnregisteredUnindexedTag construct and returns UnregisteredUnindexedTag
func NewUnregisteredUnindexedTag(name string) (*UnregisteredUnindexedTag, error) {
	if name == "" {
		return nil, fmt.Errorf("failed to create tag: name is empty")
	}
	return &UnregisteredUnindexedTag{Name: name}, nil
}

func (t *UnindexedTag) GetID() uint64 {
	return uint64(t.ID)
}

func (t *UnindexedTag) SetID(id uint64) {
	t.ID = TagID(id)
}

type Tag struct {
	*UnindexedTag
	Index int `json:"index"`
}

// NewTag construct and return Tag
func NewTag(id TagID, name string, index int) (*Tag, error) {
	errMsg := "failed to create Tag"
	if index < 0 {
		return nil, fmt.Errorf("%s: negative index is not allowed(%d)", errMsg, index)
	}
	tag, err := NewUnindexedTag(id, name)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errMsg, err)
	}
	return &Tag{UnindexedTag: tag, Index: index}, nil
}

func (t *Tag) Unregister() *UnregisteredTag {
	return &UnregisteredTag{
		UnregisteredUnindexedTag: t.UnindexedTag.Unregister(),
		Index:                    t.Index,
	}
}

func (t *Tag) SafeUnregister() (*UnregisteredTag, error) {
	if t.ID != 0 {
		return nil, fmt.Errorf("failed to unregister tag because it has ID(%d)", t.ID)
	}
	return t.Unregister(), nil
}

func (t *Tag) ReRegister(id TagID) *Tag {
	newTag := *t
	newTag.ID = id
	return &newTag
}

func (t *Tag) Unindex() *UnindexedTag {
	return t.UnindexedTag
}

type UnregisteredTag struct {
	*UnregisteredUnindexedTag
	Index int
}

func (t *UnregisteredTag) Register(id TagID) *Tag {
	return &Tag{
		UnindexedTag: t.UnregisteredUnindexedTag.Register(id),
		Index:        t.Index,
	}
}

// NewUnregisteredTag construct and return UnregisteredTag
func NewUnregisteredTag(name string, index int) (*UnregisteredTag, error) {
	errMsg := "failed to create Tag"
	if index < 0 {
		return nil, fmt.Errorf("%s: negative index is not allowed(%d)", errMsg, index)
	}
	tag, err := NewUnregisteredUnindexedTag(name)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errMsg, err)
	}
	return &UnregisteredTag{UnregisteredUnindexedTag: tag, Index: index}, nil
}

// NewUnregisteredTagWithIndexFromUnregisteredTag construct and return UnregisteredTag from UnregisteredUnindexedTag
func NewUnregisteredTagWithIndexFromUnregisteredTag(tag *UnregisteredUnindexedTag, index int) (*UnregisteredTag, error) {
	errMsg := "failed to create Tag"
	if index < 0 {
		return nil, fmt.Errorf("%s: negative index is not allowed(%d)", errMsg, index)
	}
	return &UnregisteredTag{UnregisteredUnindexedTag: tag, Index: index}, nil
}

type TagSet struct {
	m     map[TagID]*Tag
	nameM map[string]*Tag
}

// NewTagSet returns TagSet.
// if nil is provided as tags, empty TagSet will be return.
func NewTagSet(tags []*Tag) *TagSet {
	tagSet := &TagSet{
		m:     map[TagID]*Tag{},
		nameM: map[string]*Tag{},
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
func (t *TagSet) Set(tag *Tag) bool {
	if sameNamedTag, ok := t.nameM[tag.Name]; ok && sameNamedTag.ID != tag.ID {
		return false
	}
	t.m[tag.ID] = tag
	t.nameM[tag.Name] = tag
	return true
}

func (t *TagSet) Get(id TagID) (*Tag, bool) {
	tag, ok := t.m[id]
	return tag, ok
}

func (t *TagSet) GetByName(name string) (*Tag, bool) {
	tag, ok := t.nameM[name]
	return tag, ok
}

func (t *TagSet) SubSetBy(f func(tag *Tag) bool) *TagSet {
	subset := NewTagSet(nil)
	for _, tag := range t.m {
		if f(tag) {
			subset.Set(tag)
		}
	}
	return subset
}

// SplitBy splits TagSet to two sub sets based on return value of provided function.
func (t *TagSet) SplitBy(f func(tag *Tag) bool) (trueTagSet, falseTagSet *TagSet) {
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
	return t.SubSetBy(func(tag *Tag) bool {
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
	return t.SubSetBy(func(tag *Tag) bool {
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
	return t.SplitBy(func(tag *Tag) bool {
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
	return t.SplitBy(func(tag *Tag) bool {
		_, ok := m[tag.ID]
		return ok
	})
}

// ToMap returns two maps.
func (t *TagSet) ToMap() (map[TagID]*Tag, map[string]*Tag) {
	return t.m, t.nameM
}

// ToTags returns tags
func (t *TagSet) ToTags() (tags []*Tag) {
	for _, tag := range t.m {
		tags = append(tags, tag)
	}
	return
}
