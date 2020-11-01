package action

import (
	"fmt"
	"log"

	"github.com/mitchellh/mapstructure"

	"github.com/mpppk/imagine/domain/repository"

	"github.com/mpppk/imagine/domain/model"

	fsa "github.com/mpppk/lorca-fsa/lorca-fsa"
)

const (
	workSpacePrefix                     = "WORKSPACE/"
	WorkSpaceScanRequestType   fsa.Type = workSpacePrefix + "SCAN/REQUEST"
	WorkSpaceSelectSpaceType   fsa.Type = workSpacePrefix + "SELECT"
	WorkSpaceScanResultType    fsa.Type = workSpacePrefix + "SCAN/RESULT"
	WorkSpaceUpdateRequestType fsa.Type = workSpacePrefix + "UPDATE/REQUEST"
	WorkSpaceUpdateType        fsa.Type = workSpacePrefix + "UPDATE"
)

const defaultWorkSpaceName = "default-workspace"

type wsPayload struct {
	WorkSpaceName model.WSName `json:"workSpaceName"`
}

type WorkSpaceUpdatePayload struct {
	WorkSpace model.WorkSpace
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

func (w *workspaceActionCreator) Update(workSpace *model.WorkSpace) *fsa.Action {
	return &fsa.Action{
		Type:    WorkSpaceUpdateType,
		Payload: workSpace,
	}
}

type workspaceScanHandler struct {
	globalRepository repository.WorkSpace
	action           *workspaceActionCreator
}

func (d *workspaceScanHandler) Do(action *fsa.Action, dispatch fsa.Dispatch) error {
	workspaces, err := d.globalRepository.List()
	if err != nil {
		return fmt.Errorf("failed to fetch workspaces from repository: %w", err)
	}

	if len(workspaces) == 0 {
		log.Println("debug: default workspace will be created because no workspace exists")
		defaultWorkSpace := &model.WorkSpace{Name: defaultWorkSpaceName}
		if err := d.globalRepository.Add(defaultWorkSpace); err != nil {
			return fmt.Errorf("failed to create default workspace: %w", err)
		}
		workspaces = append(workspaces, defaultWorkSpace)
	}

	return dispatch(d.action.ScanResult(workspaces))
}

type workspaceUpdateHandler struct {
	globalRepository repository.WorkSpace
	action           *workspaceActionCreator
}

func (d *workspaceUpdateHandler) Do(action *fsa.Action, dispatch fsa.Dispatch) error {
	var payload WorkSpaceUpdatePayload
	if err := mapstructure.Decode(action.Payload, &payload); err != nil {
		return fmt.Errorf("failed to decode payload: %w", err)
	}

	if err := d.globalRepository.Update(&payload.WorkSpace); err != nil {
		return fmt.Errorf("failed to update workspace from repository: %w", err)
	}

	return dispatch(d.action.Update(&payload.WorkSpace))
}

type workspaceHandlerCreator struct {
	globalRepository repository.WorkSpace
	action           *workspaceActionCreator
}

func newWorkspaceHandlerCreator(globalRepository repository.WorkSpace) *workspaceHandlerCreator {
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

func (w *workspaceHandlerCreator) Update() *workspaceUpdateHandler {
	return &workspaceUpdateHandler{
		globalRepository: w.globalRepository,
		action:           w.action,
	}
}
