package action

import (
	"github.com/mpppk/imagine/domain/model"
	fsa "github.com/mpppk/lorca-fsa/lorca-fsa"
)

const (
	IndexPrefix                               = "INDEX/"
	IndexClickAddDirectoryButtonType fsa.Type = IndexPrefix + "CLICK_ADD_DIRECTORY_BUTTON"
)

const (
	AssetPrefix                     = "ASSET/"
	AssetRequestAssetsType fsa.Type = AssetPrefix + "REQUEST_ASSETS"
)

const (
	GlobalPrefix                         = "GLOBAL/"
	GlobalRequestWorkSpacesType fsa.Type = GlobalPrefix + "REQUEST_WORKSPACES"
)

const (
	ServerPrefix                               = "SERVER/"
	ServerStartDirectoryScanningType  fsa.Type = ServerPrefix + "START_DIRECTORY_SCANNING"
	ServerCancelDirectoryScanningType fsa.Type = ServerPrefix + "CANCEL_DIRECTORY_SCANNING"
	ServerFinishDirectoryScanningType fsa.Type = ServerPrefix + "FINISH_DIRECTORY_SCANNING"
	ServerScanningImagesType          fsa.Type = ServerPrefix + "SCANNING_IMAGES"
	ServerScanWorkSpaces              fsa.Type = ServerPrefix + "SCAN_WORKSPACES"
	ServerScanningAssetsType          fsa.Type = ServerPrefix + "SCANNING_ASSETS"
	ServerFinishAssetsScanningType    fsa.Type = ServerPrefix + "FINISH_ASSETS_SCANNING"
)

type WSPayload struct {
	WorkSpaceName model.WSName `json:"workSpaceName"`
}

func newWSPayload(name model.WSName) *WSPayload {
	return &WSPayload{WorkSpaceName: name}
}
