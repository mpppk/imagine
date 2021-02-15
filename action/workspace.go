package action

import (
	"fmt"
	"log"

	"github.com/mpppk/imagine/usecase/interactor"

	"github.com/mitchellh/mapstructure"

	"github.com/mpppk/imagine/domain/repository"

	"github.com/mpppk/imagine/domain/model"

	fsa "github.com/mpppk/lorca-fsa"
)

const (
	workSpacePrefix                     = "WORKSPACE/"
	WorkSpaceScanRequestType   fsa.Type = workSpacePrefix + "SCAN/REQUEST"
	WorkSpaceSelectType        fsa.Type = workSpacePrefix + "SELECT"
	WorkSpaceScanResultType    fsa.Type = workSpacePrefix + "SCAN/RESULT"
	WorkSpaceUpdateRequestType fsa.Type = workSpacePrefix + "UPDATE/REQUEST"
	WorkSpaceUpdateType        fsa.Type = workSpacePrefix + "UPDATE"
)

const defaultWorkSpaceName = "default-workspace"

type WsPayload struct {
	WorkSpaceName model.WSName `json:"workSpaceName"`
}

type WorkSpaceScanResultPayload struct {
	BasePath   string             `json:"basePath"`
	WorkSpaces []*model.WorkSpace `json:"workspaces"`
}

type WorkSpaceUpdatePayload struct {
	WorkSpace model.WorkSpace
}

type workspaceActionCreator struct{}

func newWSPayload(name model.WSName) *WsPayload {
	return &WsPayload{WorkSpaceName: name}
}

func (w *workspaceActionCreator) ScanResult(workSpaces []*model.WorkSpace) *fsa.Action {
	return &fsa.Action{
		Type: WorkSpaceScanResultType,
		Payload: &WorkSpaceScanResultPayload{
			BasePath:   interactor.DefaultBasePath,
			WorkSpaces: workSpaces,
		},
	}
}

func (w *workspaceActionCreator) Update(workSpace *model.WorkSpace) *fsa.Action {
	return &fsa.Action{
		Type:    WorkSpaceUpdateType,
		Payload: workSpace,
	}
}

type workspaceScanHandler struct {
	client *repository.Client
	action *workspaceActionCreator
}

func (d *workspaceScanHandler) Do(action *fsa.Action, dispatch fsa.Dispatch) error {
	workspaces, err := d.client.WorkSpace.List()
	if err != nil {
		return fmt.Errorf("failed to fetch workspaces from repository: %w", err)
	}

	if len(workspaces) == 0 {
		log.Println("debug: default workspace will be created because no workspace exists")
		defaultWorkSpace, err := d.client.CreateWorkSpace(defaultWorkSpaceName)
		if err != nil {
			return fmt.Errorf("failed to create default workspace: %w", err)
		}
		workspaces = append(workspaces, defaultWorkSpace)
	}

	return dispatch(d.action.ScanResult(workspaces))
}

type workspaceUpdateHandler struct {
	client *repository.Client
	action *workspaceActionCreator
}

func (d *workspaceUpdateHandler) Do(action *fsa.Action, dispatch fsa.Dispatch) error {
	var payload WorkSpaceUpdatePayload
	if err := mapstructure.Decode(action.Payload, &payload); err != nil {
		return fmt.Errorf("failed to decode payload: %w", err)
	}

	if err := d.client.WorkSpace.Update(&payload.WorkSpace); err != nil {
		return fmt.Errorf("failed to update workspace from repository: %w", err)
	}

	return dispatch(d.action.Update(&payload.WorkSpace))
}

type workspaceHandlerCreator struct {
	client *repository.Client
	action *workspaceActionCreator
}

func newWorkspaceHandlerCreator(client *repository.Client) *workspaceHandlerCreator {
	return &workspaceHandlerCreator{
		client: client,
		action: &workspaceActionCreator{},
	}
}

func (w *workspaceHandlerCreator) Scan() *workspaceScanHandler {
	return &workspaceScanHandler{
		client: w.client,
		action: w.action,
	}
}

func (w *workspaceHandlerCreator) Update() *workspaceUpdateHandler {
	return &workspaceUpdateHandler{
		client: w.client,
		action: w.action,
	}
}
