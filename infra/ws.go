package infra

import (
	"fmt"
	"net/http"

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

func ws(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	dispatch := CreateDispatch(ws)

	for {
		// Read
		var action Action
		if err := ws.ReadJSON(&action); err != nil {
			if websocket.IsCloseError(err) {
				c.Logger().Print("connection closed")
			} else {
				c.Logger().Error(err)
			}
			break
		}
		if err := HandleAction(action, dispatch); err != nil {
			c.Logger().Error(err)
		}

		fmt.Printf("%s\n", action)
	}
	return nil
}
