package action

import fsa "github.com/mpppk/lorca-fsa"

func NewHandlers() *fsa.Handlers {
	handlers := fsa.NewHandlers()
	handlers.Handle(indexClickAddDirectoryButtonType, fsa.HandlerFunc(readDirRequestHandler))
	return handlers
}
