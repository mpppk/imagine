package infra

import (
	"fmt"

	"github.com/sqweek/dialog"

	"github.com/gorilla/websocket"
)

type ActionType string

const (
	ClickAddDirectoryButton ActionType = "INDEX/CLICK_ADD_DIRECTORY_BUTTON"
)

const (
	StartDirectoryScanning ActionType = "SERVER/START_DIRECTORY_SCANNING"
	ScanningImages         ActionType = "SERVER/SCANNING_IMAGES"
)

type Action struct {
	Type    ActionType `json:"type"`
	Payload string     `json:"payload"`
	Error   string     `json:"error"`
	Meta    string     `json:"meta"`
}

func newStartDirectoryScanningAction() *Action {
	return &Action{Type: StartDirectoryScanning}
}

type ScanningImagesAction struct {
	Type    ActionType `json:"type"`
	Payload []string   `json:"payload"`
}

func newScanningImages(paths []string) *ScanningImagesAction {
	return &ScanningImagesAction{Type: ScanningImages, Payload: paths}
}

func HandleAction(action Action, dispatch func(action interface{}) error) error {
	switch action.Type {
	case ClickAddDirectoryButton:
		return handleClickAddDirectoryButton(action, dispatch)
		//return dispatch(newStartDirectoryScanningAction())
	}
	return nil
}

func CreateDispatch(ws *websocket.Conn) func(interface{}) error {
	return func(action interface{}) error {
		fmt.Println("action", action)
		return ws.WriteJSON(action)
	}
}

func handleClickAddDirectoryButton(action Action, dispatch func(action interface{}) error) error {
	go func() {
		directory, err := dialog.Directory().Title("Load images").Browse()
		if err != nil {
			panic(err)
		}
		var paths []string
		for p := range LoadImagesFromDir(directory, 10) {
			paths = append(paths, p)
			if len(paths) >= 20 {
				dispatch(newScanningImages(paths)) // FIXME
			}
		}
		if len(paths) > 0 {
			dispatch(newScanningImages(paths)) // FIXME
		}
	}()
	//ok := dialog.Message("%s", "Do you want to continue?").Title("Are you sure?").YesNo()
	return dispatch(newStartDirectoryScanningAction())
}
