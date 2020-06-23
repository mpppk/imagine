package actions

import (
	"github.com/mpppk/imagine/usecase"
	"github.com/mpppk/imagine/util"
	"github.com/sqweek/dialog"
)

type Dispatch func(action interface{}) error

type Handler struct {
	assetUseCase *usecase.Asset
	dispatch     Dispatch
}

func NewActionHandler(assetUseCase *usecase.Asset, dispatch Dispatch) *Handler {
	return &Handler{
		assetUseCase: assetUseCase,
		dispatch:     dispatch,
	}
}

func (a *Handler) clickAddDirectoryButton() error {
	if err := a.dispatch(newStartDirectoryScanningAction()); err != nil {
		return err
	}
	directory, err := dialog.Directory().Title("Load images").Browse()
	if err != nil {
		return err
	}
	var paths []string
	for p := range util.LoadImagesFromDir(directory, 10) {
		if err := a.assetUseCase.AddImage(p); err != nil {
			return err
		}
		paths = append(paths, p)
		if len(paths) >= 20 {
			if err := a.dispatch(newScanningImages(paths)); err != nil {
				return err
			}
		}
	}
	if len(paths) > 0 {
		if err := a.dispatch(newScanningImages(paths)); err != nil {
			return err
		}
	}
	return nil
}

func (a *Handler) Handle(action Action) error {
	switch action.Type {
	case ClickAddDirectoryButton:
		return a.clickAddDirectoryButton()
	}
	return nil
}

//func handleClickAddDirectoryButton(action Action, dispatch func(action interface{}) error) error {
//	go func() {
//		directory, err := dialog.Directory().Title("Load images").Browse()
//		if err != nil {
//			panic(err)
//		}
//		var paths []string
//		for p := range infra.LoadImagesFromDir(directory, 10) {
//			paths = append(paths, p)
//			if len(paths) >= 20 {
//				dispatch(newScanningImages(paths)) // FIXME
//			}
//		}
//		if len(paths) > 0 {
//			dispatch(newScanningImages(paths)) // FIXME
//		}
//	}()
//	//ok := dialog.Message("%s", "Do you want to continue?").Title("Are you sure?").YesNo()
//	return dispatch(newStartDirectoryScanningAction())
//}
