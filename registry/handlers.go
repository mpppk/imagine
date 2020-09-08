//+build !wireinject

package registry

import (
	"github.com/mpppk/imagine/action"
	fsa "github.com/mpppk/lorca-fsa/lorca-fsa"
	"go.etcd.io/bbolt"
)

func NewHandlers(db *bbolt.DB) *fsa.Handlers {
	handlerCreator := InitializeHandlerCreator(db)
	handlers := fsa.NewHandlers()
	handlers.Handle(action.IndexAddDirectoryButtonClickType, handlerCreator.FS.Scan())
	handlers.Handle(action.WorkSpaceScanRequestType, handlerCreator.Workspace.Scan())
	handlers.Handle(action.AssetScanRequestType, handlerCreator.Asset.Scan())
	handlers.Handle(action.WorkSpaceSelectSpaceType, handlerCreator.Tag.Scan())
	handlers.Handle(action.TagUpdateType, handlerCreator.Tag.Save())
	handlers.Handle(action.BoxAssignRequestType, handlerCreator.Box.Assign())
	handlers.Handle(action.BoxUnAssignRequestType, handlerCreator.Box.UnAssign())
	return handlers
}
