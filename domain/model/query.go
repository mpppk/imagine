package model

import (
	"log"
	"strings"
)

type QueryOP string

const (
	EqualsQueryOP     QueryOP = "equals"
	NotEqualsQueryOP  QueryOP = "not-equals"
	StartWithQueryOP  QueryOP = "start-with"
	NoTagsQueryOP     QueryOP = "no-tags"
	PathEqualsQueryOP QueryOP = "path-equals"
)

type Query struct {
	Op    QueryOP `json:"op"`
	Value string  `json:"value"`
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
	case NoTagsQueryOP:
		return len(asset.BoundingBoxes) == 0
	case PathEqualsQueryOP:
		return true // FIXME
	default:
		log.Printf("warning: unknown query op is given: %s", q.Op)
		return false
	}
}
