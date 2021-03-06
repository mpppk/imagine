package option

import (
	"fmt"
	"strings"

	"github.com/mpppk/imagine/domain/model"
)

// parseQuery parse query string like "start-with=tag1,equals=tag2"
func parseQuery(query string) (queries []*model.Query, err error) {
	errMsg := "failed to parse query"
	for _, query2 := range strings.Split(query, ",") {
		if query2 == "" { // for trailing comma
			continue
		}
		opAndValue := strings.Split(query2, "=")
		if len(opAndValue) != 2 {
			return nil, fmt.Errorf("invalid query format: %s", opAndValue)
		}
		op, value := opAndValue[0], opAndValue[1]
		query, err := model.NewQuery(op, value)
		if err != nil {
			return nil, fmt.Errorf("%s, %w", errMsg, err)
		}
		queries = append(queries, query)
	}
	return
}
