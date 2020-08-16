package action

import (
	"fmt"
	"log"

	"github.com/mpppk/imagine/domain/repository"

	"github.com/mpppk/imagine/domain/model"

	fsa "github.com/mpppk/lorca-fsa/lorca-fsa"
)

const (
	workSpacePrefix                   = "WORKSPACE/"
	WorkSpaceScanRequestType fsa.Type = workSpacePrefix + "SCAN/REQUEST"
	WorkSpaceSelectSpaceType fsa.Type = workSpacePrefix + "SELECT"
	WorkSpaceScanResultType  fsa.Type = workSpacePrefix + "SCAN/RESULT"
)

const defaultWorkSpaceName = "default-workspace"

type wsPayload struct {
	WorkSpaceName model.WSName `json:"workSpaceName"`
}

type workspaceActionCreator struct{}

func newWSPayload(name model.WSName) *wsPayload {
	return &wsPayload{WorkSpaceName: name}
}

func (w *workspaceActionCreator) ScanResult(workSpaces []*model.WorkSpace) *fsa.Action {
	return &fsa.Action{
		Type:    WorkSpaceScanResultType,
		Payload: workSpaces,
	}
}

type workspaceScanHandler struct {
	globalRepository repository.Global
	action           *workspaceActionCreator
}

func (d *workspaceScanHandler) Do(action *fsa.Action, dispatch fsa.Dispatch) error {
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

	return dispatch(d.action.ScanResult(workspaces))
}

type workspaceHandlerCreator struct {
	globalRepository repository.Global
	action           *workspaceActionCreator
}

func newWorkspaceHandlerCreator(globalRepository repository.Global) *workspaceHandlerCreator {
	return &workspaceHandlerCreator{
		globalRepository: globalRepository,
		action:           &workspaceActionCreator{},
	}
}

func (w *workspaceHandlerCreator) Scan() *workspaceScanHandler {
	return &workspaceScanHandler{
		globalRepository: w.globalRepository,
		action:           w.action,
	}
}
