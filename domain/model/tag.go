package model

type TagID uint64
type Tag struct {
	ID   TagID  `json:"id"`
	Name string `json:"name"`
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

// SubSetByNames returns sub TagSet. new TagSet contains tag which have either provided names.
func (t *TagSet) SubSetByNames(names []string) *TagSet {
	return t.SubSetBy(func(tag *Tag) bool {
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
	return t.SplitBy(func(tag *Tag) bool {
		for _, name := range names {
			if tag.Name == name {
				return true
			}
		}
		return false
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
