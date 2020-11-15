package repository

import "github.com/mpppk/imagine/domain/model"

type WorkSpace interface {
	Init() error
	List() (workspaces []*model.WorkSpace, err error)
	Add(ws *model.WorkSpace) error
	Update(ws *model.WorkSpace) error
}
