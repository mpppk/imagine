package handler

import (
	"fmt"
	"net/http"

	"github.com/mpppk/imagine/actions"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		// FIXME need for debug
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func (h *Handlers) WS(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	dispatch := actions.CreateDispatch(ws)
	actionHandler := actions.NewActionHandler(h.assetUseCase, dispatch)

	for {
		var action actions.Action
		if err := ws.ReadJSON(&action); err != nil {
			if websocket.IsCloseError(err) {
				c.Logger().Print("connection closed")
			} else {
				c.Logger().Error(err)
			}
			break
		}
		if err := actionHandler.Handle(action); err != nil {
			c.Logger().Error(err)
		}

		fmt.Printf("%s\n", action)
	}
	return nil
}
