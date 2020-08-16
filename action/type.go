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
	WorkSpacePrefix                         = "WORKSPACE/"
	WorkSpaceRequestWorkSpacesType fsa.Type = WorkSpacePrefix + "REQUEST_WORKSPACES"
	WorkSpaceSelectNewWorkSpace    fsa.Type = WorkSpacePrefix + "SELECT_NEW_WORKSPACE"
)

const (
	FSPrefix                   = "FS/"
	FSScanCancelType  fsa.Type = FSPrefix + "SCAN/CANCEL"
	FSScanStartType   fsa.Type = FSPrefix + "SCAN/START"
	FSScanFinishType  fsa.Type = FSPrefix + "SCAN/FINISH"
	FSScanRunningType fsa.Type = FSPrefix + "SCAN/RUNNING"
)

const (
	ServerPrefix = "SERVER/"
	//ServerScanningImagesType fsa.Type = ServerPrefix + "SCANNING_IMAGES"
	ServerScanWorkSpaces     fsa.Type = ServerPrefix + "SCAN_WORKSPACES"
	ServerScanningAssetsType fsa.Type = ServerPrefix + "SCANNING_ASSETS"
	ServerTagScanType        fsa.Type = ServerPrefix + "TAG/SCAN"
	ServerTagSaveType        fsa.Type = ServerPrefix + "TAG/SAVE"
)

type WSPayload struct {
	WorkSpaceName model.WSName `json:"workSpaceName"`
}

func newWSPayload(name model.WSName) *WSPayload {
	return &WSPayload{WorkSpaceName: name}
}
