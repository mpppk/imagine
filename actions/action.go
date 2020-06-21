package actions

import (
	"fmt"

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

func CreateDispatch(ws *websocket.Conn) Dispatch {
	return func(action interface{}) error {
		fmt.Println("action", action)
		return ws.WriteJSON(action)
	}
}
