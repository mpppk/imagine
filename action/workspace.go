package action

import (
	"fmt"
	"log"

	"github.com/mpppk/imagine/domain/repository"

	"github.com/mpppk/imagine/domain/model"

	fsa "github.com/mpppk/lorca-fsa/lorca-fsa"
)

const defaultWorkSpaceName = "default-workspace"

func newScanWorkSpacesAction(workSpaces []*model.WorkSpace) *fsa.Action {
	return &fsa.Action{
		Type:    ServerScanWorkSpaces,
		Payload: workSpaces,
	}
}

type RequestWorkSpacesHandler struct {
	globalRepository repository.Global
}

func NewRequestWorkSpacesHandler(globalRepository repository.Global) *RequestWorkSpacesHandler {
	return &RequestWorkSpacesHandler{globalRepository: globalRepository}
}

func (d *RequestWorkSpacesHandler) Do(action *fsa.Action, dispatch fsa.Dispatch) error {
	workspaces, err := d.globalRepository.ListWorkSpace()
	if err != nil {
		return fmt.Errorf("failed to fetch workspaces from repository: %w", err)
	}

	if len(workspaces) == 0 {
		log.Println("debug: default workspace will be created because no workspace exists")
		defaultWorkSpace := &model.WorkSpace{Name: defaultWorkSpaceName}
		if err := d.globalRepository.AddWorkSpace(defaultWorkSpace); err != nil {
			return fmt.Errorf("failed to create default workspace: %w", err)
		}
		workspaces = append(workspaces, defaultWorkSpace)
	}

	return dispatch(newScanWorkSpacesAction(workspaces))
}
