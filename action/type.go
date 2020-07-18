package action

import fsa "github.com/mpppk/lorca-fsa"

const (
	IndexPrefix                               = "INDEX/"
	IndexClickAddDirectoryButtonType fsa.Type = IndexPrefix + "CLICK_ADD_DIRECTORY_BUTTON"
)

const (
	ServerPrefix                               = "SERVER/"
	ServerStartDirectoryScanningType  fsa.Type = ServerPrefix + "START_DIRECTORY_SCANNING"
	ServerCancelDirectoryScanningType fsa.Type = ServerPrefix + "CANCEL_DIRECTORY_SCANNING"
	ServerFinishDirectoryScanningType fsa.Type = ServerPrefix + "FINISH_DIRECTORY_SCANNING"
	ServerScanningImagesType          fsa.Type = ServerPrefix + "SCANNING_IMAGES"
)
