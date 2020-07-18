package action

import (
	"fmt"

	"github.com/mpppk/imagine/usecase"

	"github.com/gen2brain/dlgs"
	"github.com/mpppk/imagine/util"
	fsa "github.com/mpppk/lorca-fsa"
)

func newStartDirectoryScanningAction() *fsa.Action {
	return &fsa.Action{
		Type: ServerStartDirectoryScanningType,
	}
}

func newFinishDirectoryScanningAction() *fsa.Action {
	return &fsa.Action{
		Type: ServerFinishDirectoryScanningType,
	}
}

func newCancelDirectoryScanningAction() *fsa.Action {
	return &fsa.Action{
		Type: ServerCancelDirectoryScanningType,
	}
}

func newScanningImages(paths []string) *fsa.Action {
	return &fsa.Action{
		Type:    ServerScanningImagesType,
		Payload: paths,
	}
}

type DirectoryScanHandler struct {
	assetUseCase *usecase.Asset
}

func NewReadDirectoryScanHandler(assetUseCase *usecase.Asset) *DirectoryScanHandler {
	return &DirectoryScanHandler{assetUseCase: assetUseCase}
}

func (d *DirectoryScanHandler) Do(action *fsa.Action, dispatch fsa.Dispatch) error {
	fmt.Println(action)

	if err := dispatch(newStartDirectoryScanningAction()); err != nil {
		return err
	}

	directory, selected, err := dlgs.File("Select file", "", true)
	if err != nil {
		return fmt.Errorf("failed to open file selector: %w", err)
	}

	if !selected {
		return dispatch(newCancelDirectoryScanningAction())
	}

	var paths []string
	for p := range util.LoadImagesFromDir(directory, 10) {
		if err := d.assetUseCase.AddImage(p); err != nil {
			return err
		}
		paths = append(paths, p)
		if len(paths) >= 20 {
			if err := dispatch(newScanningImages(paths)); err != nil {
				return err
			}
		}
	}
	if len(paths) > 0 {
		if err := dispatch(newScanningImages(paths)); err != nil {
			return err
		}
	}

	return dispatch(newFinishDirectoryScanningAction())
}
