package action

import (
	"fmt"

	"github.com/mpppk/imagine/domain/model"

	"github.com/mitchellh/mapstructure"
	"github.com/mpppk/imagine/usecase"
	fsa "github.com/mpppk/lorca-fsa/lorca-fsa"
)

const (
	assetPrefix                     = "ASSET/"
	AssetRequestAssetsType fsa.Type = assetPrefix + "REQUEST_ASSETS"
	AssetScanRunningType   fsa.Type = assetPrefix + "SCAN/RUNNING"
)

type RequestAssetsHandler struct {
	c            <-chan *model.Asset
	assetUseCase *usecase.Asset
}

func NewRequestAssetsHandler(assetUseCase *usecase.Asset) *RequestAssetsHandler {
	return &RequestAssetsHandler{assetUseCase: assetUseCase}
}

type RequestAssetPayload struct {
	wsPayload  `mapstructure:",squash"`
	RequestNum int `json:"RequestNum"`
}

type ScanningAssetsPayload struct {
	*wsPayload
	Assets []*model.Asset `json:"assets"`
}

func newAssetScanRunning(wsName model.WSName, assets []*model.Asset) *fsa.Action {
	return &fsa.Action{
		Type: AssetScanRunningType,
		Payload: &ScanningAssetsPayload{
			wsPayload: newWSPayload(wsName),
			Assets:    assets,
		},
	}
}

func newFinishAssetScanningType(wsName model.WSName) *fsa.Action {
	return &fsa.Action{
		Type:    FSScanFinishType, // FIXME
		Payload: newWSPayload(wsName),
	}
}

func (d *RequestAssetsHandler) Do(action *fsa.Action, dispatch fsa.Dispatch) error {
	var payload RequestAssetPayload
	if err := mapstructure.Decode(action.Payload, &payload); err != nil {
		return fmt.Errorf("failed to decode payload: %w", err)
	}

	if d.c == nil {
		c, err := d.assetUseCase.ListAsync(payload.WorkSpaceName)
		if err != nil {
			return err
		}
		d.c = c
	}

	var ret []*model.Asset
	for asset := range d.c {
		ret = append(ret, asset)
		if len(ret) >= payload.RequestNum {
			return dispatch(newAssetScanRunning(payload.WorkSpaceName, ret))
		}
	}
	if len(ret) > 0 {
		if err := dispatch(newAssetScanRunning(payload.WorkSpaceName, ret)); err != nil {
			return err
		}
	}
	return dispatch(newFinishAssetScanningType(payload.WorkSpaceName))
}
