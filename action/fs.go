package action

import (
	"fmt"

	"github.com/mitchellh/mapstructure"

	"github.com/mpppk/imagine/domain/model"

	"github.com/mpppk/imagine/usecase"

	"github.com/gen2brain/dlgs"
	"github.com/mpppk/imagine/util"
	fsa "github.com/mpppk/lorca-fsa/lorca-fsa"
)

const (
	fsPrefix                   = "FS/"
	FSScanCancelType  fsa.Type = fsPrefix + "SCAN/CANCEL"
	FSScanStartType   fsa.Type = fsPrefix + "SCAN/START"
	FSScanFinishType  fsa.Type = fsPrefix + "SCAN/FINISH"
	FSScanRunningType fsa.Type = fsPrefix + "SCAN/RUNNING"
)

type ScanningImagesPayload struct {
	*wsPayload
	Paths []string `json:"paths"`
}

type FsScanStartPayload struct {
	*wsPayload
	BasePath string `json:"basePath"`
}

type fsActionCreator struct{}

func (f *fsActionCreator) scanStart(wsName model.WSName, basePath string) *fsa.Action {
	return &fsa.Action{
		Type: FSScanStartType,
		Payload: FsScanStartPayload{
			wsPayload: newWSPayload(wsName),
			BasePath:  basePath,
		},
	}
}

func (f *fsActionCreator) scanFinish(wsName model.WSName) *fsa.Action {
	return &fsa.Action{
		Type:    FSScanFinishType,
		Payload: newWSPayload(wsName),
	}
}

func (f *fsActionCreator) scanCancel(wsName model.WSName) *fsa.Action {
	return &fsa.Action{
		Type:    FSScanCancelType,
		Payload: newWSPayload(wsName),
	}
}

func (f *fsActionCreator) scanRunning(wsName model.WSName, paths []string) *fsa.Action {
	return &fsa.Action{
		Type: FSScanRunningType,
		Payload: &ScanningImagesPayload{
			wsPayload: newWSPayload(wsName),
			Paths:     paths,
		},
	}
}

type fsScanHandler struct {
	assetUseCase *usecase.Asset
	action       *fsActionCreator
}

func (f *fsScanHandler) Do(action *fsa.Action, dispatch fsa.Dispatch) error {
	var payload wsPayload
	if err := mapstructure.Decode(action.Payload, &payload); err != nil {
		return fmt.Errorf("failed to decode payload: %w", err)
	}

	directory, selected, err := dlgs.File("Select file", "", true)
	if err != nil {
		return fmt.Errorf("failed to open file selector: %w", err)
	}

	if !selected {
		return dispatch(f.action.scanCancel(payload.WorkSpaceName))
	}

	if err := dispatch(f.action.scanStart(payload.WorkSpaceName, directory)); err != nil {
		return err
	}

	var paths []string
	for p := range util.LoadImagesFromDir(directory, 10) {
		if added, err := f.assetUseCase.AddAssetFromImagePathIfDoesNotExist(payload.WorkSpaceName, p); err != nil {
			return err
		} else if !added {
			continue
		}
		paths = append(paths, p)
		if len(paths) >= 20 {
			if err := dispatch(f.action.scanRunning(payload.WorkSpaceName, paths)); err != nil {
				return err
			}
			paths = []string{}
		}
	}
	if len(paths) > 0 {
		if err := dispatch(f.action.scanRunning(payload.WorkSpaceName, paths)); err != nil {
			return err
		}
	}

	return dispatch(f.action.scanFinish(payload.WorkSpaceName))
}

type fsHandlerCreator struct {
	assetUseCase *usecase.Asset
	action       *fsActionCreator
}

func newFSHandlerCreator(assetUseCase *usecase.Asset) *fsHandlerCreator {
	return &fsHandlerCreator{
		assetUseCase: assetUseCase,
		action:       &fsActionCreator{},
	}
}

func (f *fsHandlerCreator) Scan() *fsScanHandler {
	return &fsScanHandler{
		assetUseCase: f.assetUseCase,
		action:       f.action,
	}
}
