package action

import (
	"fmt"
	"log"
	"path/filepath"

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
	FSScanFailType    fsa.Type = fsPrefix + "SCAN/FAIL"
)

type ScanningImagesPayload struct {
	*wsPayload
	Paths []string `json:"paths"`
}

type FsScanStartPayload struct {
	*wsPayload
	BasePath string `json:"basePath"`
}

type FsScanFailPayload struct {
	*wsPayload
	err error
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

func (f *fsActionCreator) scanFail(wsName model.WSName, err error) *fsa.Action {
	return &fsa.Action{
		Type: FSScanFailType,
		Payload: &FsScanFailPayload{
			wsPayload: newWSPayload(wsName),
			err:       err,
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

	dispatchScanFailActionAndLogOrPanic := func(err error) {
		if e := dispatch(f.action.scanFail(payload.WorkSpaceName, err)); e != nil {
			panic(e)
		}
		log.Printf("warning: %s", err)
	}

	go func() {
		var paths []string
		cnt := 0
		for p := range util.LoadImagesFromDir(directory, 500) {
			cnt++
			relP, err := util.ToRelPath(directory, p)
			if err != nil {
				dispatchScanFailActionAndLogOrPanic(err)
				continue
			}

			paths = append(paths, filepath.Clean(relP))
			if len(paths) >= 1000 {
				if _, err := f.assetUseCase.AddAssetFromImagePathListIfDoesNotExist(payload.WorkSpaceName, paths); err != nil {
					dispatchScanFailActionAndLogOrPanic(err)
					continue
				}

				if err := dispatch(f.action.scanRunning(payload.WorkSpaceName, paths)); err != nil {
					dispatchScanFailActionAndLogOrPanic(err)
					continue
				}
				paths = []string{}
			}
		}
		if len(paths) > 0 {
			if err := dispatch(f.action.scanRunning(payload.WorkSpaceName, paths)); err != nil {
				dispatchScanFailActionAndLogOrPanic(err)
			}
		}

		if err := dispatch(f.action.scanFinish(payload.WorkSpaceName)); err != nil {
			dispatchScanFailActionAndLogOrPanic(err)
		}
	}()
	return nil
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
