package model

import (
	"fmt"
	"log"
	"strings"
)

type QueryOP string

const (
	EqualsQueryOP       QueryOP = "equals"
	NotEqualsQueryOP    QueryOP = "not-equals"
	StartWithQueryOP    QueryOP = "start-with"
	NotStartWithQueryOP QueryOP = "not-start-with"
	NoTagsQueryOP       QueryOP = "no-tags"
	PathEqualsQueryOP   QueryOP = "path-equals"
)

var ops = []QueryOP{EqualsQueryOP, NotEqualsQueryOP, StartWithQueryOP, NotStartWithQueryOP, NoTagsQueryOP, PathEqualsQueryOP}

type Query struct {
	Op    QueryOP `json:"op"`
	Value string  `json:"value"`
}

func toQueryOP(opStr string) (QueryOP, bool) {
	for _, op := range ops {
		if string(op) == opStr {
			return op, true
		}
	}
	return "", false
}

func NewQuery(opStr, value string) (*Query, error) {
	op, ok := toQueryOP(opStr)
	if !ok {
		return nil, fmt.Errorf("failed to create query. invalid query opStr is provided: %s", opStr)
	}
	return &Query{Op: op, Value: value}, nil
}

func (q *Query) Match(asset *Asset, tagSet *TagSet) bool {
	switch q.Op {
	case EqualsQueryOP:
		tag, ok := tagSet.GetByName(q.Value)
		if !ok {
			return false
		}
		return asset.HasTag(tag.ID)
	case NotEqualsQueryOP:
		tag, ok := tagSet.GetByName(q.Value)
		if !ok {
			return false
		}
		return !asset.HasTag(tag.ID)
	case StartWithQueryOP:
		f := func(tag *Tag) bool {
			return strings.HasPrefix(tag.Name, q.Value)
		}
		return asset.HasAnyOneOfTagID(tagSet.SubSetBy(f))
	case NotStartWithQueryOP:
		f := func(tag *Tag) bool {
			return !strings.HasPrefix(tag.Name, q.Value)
		}
		return asset.HasAnyOneOfTagID(tagSet.SubSetBy(f))
	case NoTagsQueryOP:
		return len(asset.BoundingBoxes) == 0
	case PathEqualsQueryOP:
		return asset != nil && asset.Path == q.Value
	default:
		log.Printf("warning: unknown query op is given: %s", q.Op)
		return false
	}
}
