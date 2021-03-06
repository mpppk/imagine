package option

import (
	"reflect"
	"testing"

	"github.com/mpppk/imagine/domain/model"
)

func Test_parseQuery(t *testing.T) {
	type args struct {
		query string
	}
	tests := []struct {
		name        string
		args        args
		wantQueries []*model.Query
		wantErr     bool
	}{
		{
			name:        "parse 1 query",
			args:        args{query: "equals=tag1"},
			wantQueries: []*model.Query{{Op: "equals", Value: "tag1"}},
		},
		{
			name: "parse all queries",
			args: args{
				query: "equals=tag1,not-equals=tag2,start-with=tag3,no-tags=tag4,path-equals=tag5",
			},
			wantQueries: []*model.Query{
				{Op: "equals", Value: "tag1"},
				{Op: "not-equals", Value: "tag2"},
				{Op: "start-with", Value: "tag3"},
				{Op: "no-tags", Value: "tag4"},
				{Op: "path-equals", Value: "tag5"},
			},
		},
		{
			name: "allow trailing comma",
			args: args{
				query: "equals=tag1,",
			},
			wantQueries: []*model.Query{{Op: "equals", Value: "tag1"}},
		},
		{
			name: "missing equal",
			args: args{
				query: "equals,",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotQueries, err := parseQuery(tt.args.query)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseQuery() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotQueries, tt.wantQueries) {
				t.Errorf("parseQuery() gotQueries = %v, want %v", gotQueries, tt.wantQueries)
			}
		})
	}
}
