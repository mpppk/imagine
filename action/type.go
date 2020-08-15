package action

import (
	"github.com/mpppk/imagine/domain/model"
	fsa "github.com/mpppk/lorca-fsa/lorca-fsa"
)

const (
	IndexPrefix                               = "INDEX/"
	IndexClickAddDirectoryButtonType fsa.Type = IndexPrefix + "CLICK_ADD_DIRECTORY_BUTTON"
	IndexUpdateTags                  fsa.Type = IndexPrefix + "UPDATE_TAGS"
)

const (
	AssetPrefix                     = "ASSET/"
	AssetRequestAssetsType fsa.Type = AssetPrefix + "REQUEST_ASSETS"
)

const (
	TagPrefix               = "TAG/"
	TagRequestType fsa.Type = TagPrefix + "REQUEST"
)

const (
	GlobalPrefix                         = "GLOBAL/"
	GlobalRequestWorkSpacesType fsa.Type = GlobalPrefix + "REQUEST_WORKSPACES"
	GlobalSelectNewWorkSpace    fsa.Type = GlobalPrefix + "SELECT_NEW_WORKSPACE"
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
	ServerTagScanType                 fsa.Type = ServerPrefix + "TAG/SCAN"
	ServerTagSaveType                 fsa.Type = ServerPrefix + "TAG/SAVE"
)

type WSPayload struct {
	WorkSpaceName model.WSName `json:"workSpaceName"`
}

func newWSPayload(name model.WSName) *WSPayload {
	return &WSPayload{WorkSpaceName: name}
}
