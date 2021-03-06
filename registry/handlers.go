//+build !wireinject

package registry

import (
	"io"

	"github.com/mpppk/imagine/action"
	"github.com/mpppk/imagine/infra"
	fsa "github.com/mpppk/lorca-fsa"
	"go.etcd.io/bbolt"
)

func NewHandlers(db *bbolt.DB) *fsa.Handlers {
	handlerCreator := NewBoltHandlerCreator(db)
	handlers := fsa.NewHandlers()
	handlers.Handle(action.FSScanRequestType, handlerCreator.FS.Scan())
	handlers.Handle(action.WorkSpaceScanRequestType, handlerCreator.Workspace.Scan())
	handlers.Handle(action.AssetScanRequestType, handlerCreator.Asset.Scan())
	handlers.Handle(action.WorkSpaceSelectType, handlerCreator.Tag.Scan())
	handlers.Handle(action.FSBaseDirSelectType, handlerCreator.FS.Serve())
	handlers.Handle(action.IndexChangeBasePathButtonClickType, handlerCreator.FS.BaseDirDialog())
	handlers.Handle(action.WorkSpaceUpdateRequestType, handlerCreator.Workspace.Update())
	handlers.Handle(action.TagUpdateType, handlerCreator.Tag.Save())
	handlers.Handle(action.BoxAssignRequestType, handlerCreator.Box.Assign())
	handlers.Handle(action.BoxUnAssignRequestType, handlerCreator.Box.UnAssign())
	handlers.Handle(action.BoxModifyRequestType, handlerCreator.Box.Modify())
	handlers.Handle(action.BoxDeleteRequestType, handlerCreator.Box.Delete())
	return handlers
}
func NewHandlersWithDBPath(dbPath string) (*fsa.Handlers, io.Closer, error) {
	db, err := infra.NewBoltDB(dbPath)
	if err != nil {
		return nil, nil, err
	}
	return NewHandlers(db), db, nil
}
