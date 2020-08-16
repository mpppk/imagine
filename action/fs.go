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

func newFSScanStartAction(wsName model.WSName) *fsa.Action {
	return &fsa.Action{
		Type:    FSScanStartType,
		Payload: newWSPayload(wsName),
	}
}

func newFSScanFinishAction(wsName model.WSName) *fsa.Action {
	return &fsa.Action{
		Type:    FSScanFinishType,
		Payload: newWSPayload(wsName),
	}
}

func newFSScanCancelAction(wsName model.WSName) *fsa.Action {
	return &fsa.Action{
		Type:    FSScanCancelType,
		Payload: newWSPayload(wsName),
	}
}

type ScanningImagesPayload struct {
	*WSPayload
	Paths []string `json:"paths"`
}

func newScanningImages(wsName model.WSName, paths []string) *fsa.Action {
	return &fsa.Action{
		Type: ServerScanningImagesType,
		Payload: &ScanningImagesPayload{
			WSPayload: newWSPayload(wsName),
			Paths:     paths,
		},
	}
}

type FSScanHandler struct {
	assetUseCase *usecase.Asset
}

func NewFSScanHandler(assetUseCase *usecase.Asset) *FSScanHandler {
	return &FSScanHandler{assetUseCase: assetUseCase}
}

func (d *FSScanHandler) Do(action *fsa.Action, dispatch fsa.Dispatch) error {
	var payload WSPayload
	if err := mapstructure.Decode(action.Payload, &payload); err != nil {
		return fmt.Errorf("failed to decode payload: %w", err)
	}

	if err := dispatch(newFSScanStartAction(payload.WorkSpaceName)); err != nil {
		return err
	}

	directory, selected, err := dlgs.File("Select file", "", true)
	if err != nil {
		return fmt.Errorf("failed to open file selector: %w", err)
	}

	if !selected {
		return dispatch(newFSScanCancelAction(payload.WorkSpaceName))
	}

	var paths []string
	for p := range util.LoadImagesFromDir(directory, 10) {
		if err := d.assetUseCase.AddImage(payload.WorkSpaceName, p); err != nil {
			return err
		}
		paths = append(paths, p)
		if len(paths) >= 20 {
			if err := dispatch(newScanningImages(payload.WorkSpaceName, paths)); err != nil {
				return err
			}
		}
	}
	if len(paths) > 0 {
		if err := dispatch(newScanningImages(payload.WorkSpaceName, paths)); err != nil {
			return err
		}
	}

	return dispatch(newFSScanFinishAction(payload.WorkSpaceName))
}
