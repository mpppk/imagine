package model

import "log"

type QueryOP string

const (
	EqualsQueryOP    QueryOP = "equals"
	NotEqualsQueryOP QueryOP = "not-equals"
)

type Query struct {
	Op      QueryOP `json:"op"`
	TagName string  `json:"tagName"`
}

func (q *Query) Match(asset *Asset) bool {
	switch q.Op {
	case EqualsQueryOP:
		return asset.HasTag(q.TagName)
	case NotEqualsQueryOP:
		return !asset.HasTag(q.TagName)
	default:
		log.Printf("warning: unknown query op is given: %s", q.Op)
		return false
	}
}
