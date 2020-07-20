package action

import (
	"fmt"

	"github.com/mitchellh/mapstructure"

	"github.com/mpppk/imagine/domain/model"

	"github.com/mpppk/imagine/usecase"

	"github.com/gen2brain/dlgs"
	"github.com/mpppk/imagine/util"
	fsa "github.com/mpppk/lorca-fsa"
)

func newStartDirectoryScanningAction(wsName model.WSName) *fsa.Action {
	return &fsa.Action{
		Type:    ServerStartDirectoryScanningType,
		Payload: newWSPayload(wsName),
	}
}

func newFinishDirectoryScanningAction(wsName model.WSName) *fsa.Action {
	return &fsa.Action{
		Type:    ServerFinishDirectoryScanningType,
		Payload: newWSPayload(wsName),
	}
}

func newCancelDirectoryScanningAction(wsName model.WSName) *fsa.Action {
	return &fsa.Action{
		Type:    ServerCancelDirectoryScanningType,
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

type DirectoryScanHandler struct {
	assetUseCase *usecase.Asset
}

func NewReadDirectoryScanHandler(assetUseCase *usecase.Asset) *DirectoryScanHandler {
	return &DirectoryScanHandler{assetUseCase: assetUseCase}
}

func (d *DirectoryScanHandler) Do(action *fsa.Action, dispatch fsa.Dispatch) error {
	var payload WSPayload
	if err := mapstructure.Decode(action.Payload, &payload); err != nil {
		return fmt.Errorf("failed to decode payload: %w", err)
	}

	if err := dispatch(newStartDirectoryScanningAction(payload.WorkSpaceName)); err != nil {
		return err
	}

	directory, selected, err := dlgs.File("Select file", "", true)
	if err != nil {
		return fmt.Errorf("failed to open file selector: %w", err)
	}

	if !selected {
		return dispatch(newCancelDirectoryScanningAction(payload.WorkSpaceName))
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

	return dispatch(newFinishDirectoryScanningAction(payload.WorkSpaceName))
}
