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
	AssetScanRunningType   fsa.Type = AssetPrefix + "SCAN/RUNNING"
)

const (
	TagPrefix               = "TAG/"
	TagRequestType fsa.Type = TagPrefix + "REQUEST"
	TagScanType    fsa.Type = TagPrefix + "SCAN"
	TagSaveType    fsa.Type = TagPrefix + "SAVE"
)

const (
	WorkSpacePrefix                         = "WORKSPACE/"
	WorkSpaceRequestWorkSpacesType fsa.Type = WorkSpacePrefix + "REQUEST_WORKSPACES"
	WorkSpaceSelectNewWorkSpace    fsa.Type = WorkSpacePrefix + "SELECT_NEW_WORKSPACE"
	WorkSpaceScanWorkSpaces        fsa.Type = WorkSpacePrefix + "SCAN_WORKSPACES"
)

const (
	FSPrefix                   = "FS/"
	FSScanCancelType  fsa.Type = FSPrefix + "SCAN/CANCEL"
	FSScanStartType   fsa.Type = FSPrefix + "SCAN/START"
	FSScanFinishType  fsa.Type = FSPrefix + "SCAN/FINISH"
	FSScanRunningType fsa.Type = FSPrefix + "SCAN/RUNNING"
)

type WSPayload struct {
	WorkSpaceName model.WSName `json:"workSpaceName"`
}

func newWSPayload(name model.WSName) *WSPayload {
	return &WSPayload{WorkSpaceName: name}
}
