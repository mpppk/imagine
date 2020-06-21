package model

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

func HandleAction(action Action, dispatch func(action interface{}) error) error {
	switch action.Type {
	case ClickAddDirectoryButton:
		return dispatch(newStartDirectoryScanningAction())
	}
	return nil
}

func CreateDispatch(ws *websocket.Conn) func(interface{}) error {
	return func(action interface{}) error {
		fmt.Println("action", action)
		return ws.WriteJSON(action)
	}
}
