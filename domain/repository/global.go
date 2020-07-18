package repository

import "github.com/mpppk/imagine/domain/model"

type Global interface {
	ListWorkSpace() (workspaces []*model.WorkSpace, err error)
	AddWorkSpace(ws *model.WorkSpace) error
}
