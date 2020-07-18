package action

import fsa "github.com/mpppk/lorca-fsa"

const (
	indexPrefix                               = "INDEX/"
	indexClickAddDirectoryButtonType fsa.Type = indexPrefix + "CLICK_ADD_DIRECTORY_BUTTON"
)

const (
	serverPrefix                               = "SERVER/"
	serverStartDirectoryScanningType  fsa.Type = serverPrefix + "START_DIRECTORY_SCANNING"
	serverCancelDirectoryScanningType fsa.Type = serverPrefix + "CANCEL_DIRECTORY_SCANNING"
	serverFinishDirectoryScanningType fsa.Type = serverPrefix + "FINISH_DIRECTORY_SCANNING"
	serverScanningImagesType          fsa.Type = serverPrefix + "SCANNING_IMAGES"
)
