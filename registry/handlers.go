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
	handlers.Handle(action.IndexClickAddDirectoryButtonType, handlerCreator.NewFSScanHandler())
	handlers.Handle(action.WorkSpaceRequestWorkSpacesType, handlerCreator.NewRequestWorkSpacesHandler())
	handlers.Handle(action.AssetScanRequestType, handlerCreator.Asset.Scan())
	handlers.Handle(action.TagRequestType, handlerCreator.Asset.Scan())
	handlers.Handle(action.WorkSpaceSelectNewWorkSpace, handlerCreator.NewTagRequestHandler())
	handlers.Handle(action.IndexUpdateTags, handlerCreator.NewTagUpdateHandler())
	return handlers
}
