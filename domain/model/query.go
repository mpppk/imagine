package model

import (
	"log"
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

func (q *Query) Match(asset *Asset) bool {
	switch q.Op {
	case EqualsQueryOP:
		return asset.HasTag(q.Value)
	case NotEqualsQueryOP:
		return !asset.HasTag(q.Value)
	case StartWithQueryOP:
		return asset.HasTagStartWith(q.Value)
	case NoTagsQueryOP:
		return len(asset.BoundingBoxes) == 0
	case PathEqualsQueryOP:
		return true
	default:
		log.Printf("warning: unknown query op is given: %s", q.Op)
		return false
	}
}
