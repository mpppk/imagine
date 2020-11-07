package action

import (
	"context"
	"fmt"

	"github.com/mpppk/imagine/domain/model"

	"github.com/mitchellh/mapstructure"
	"github.com/mpppk/imagine/usecase"
	fsa "github.com/mpppk/lorca-fsa/lorca-fsa"
)

const (
	assetPrefix                   = "ASSET/"
	AssetScanRequestType fsa.Type = assetPrefix + "SCAN/REQUEST"
	AssetScanRunningType fsa.Type = assetPrefix + "SCAN/RUNNING"
	AssetScanFinishType  fsa.Type = assetPrefix + "SCAN/FINISH"
)

type assetScanRequestPayload struct {
	wsPayload  `mapstructure:",squash"`
	RequestNum int            `json:"RequestNum"`
	Queries    []*model.Query `json:"queries"`
	Reset      bool           `json:"reset"`
}

type assetScanRunningPayload struct {
	*wsPayload
	Assets []*model.Asset `json:"assets"`
}

type assetActionCreator struct{}

func (a *assetActionCreator) scanRunning(wsName model.WSName, assets []*model.Asset) *fsa.Action {
	return &fsa.Action{
		Type: AssetScanRunningType,
		Payload: &assetScanRunningPayload{
			wsPayload: newWSPayload(wsName),
			Assets:    assets,
		},
	}
}

func (a *assetActionCreator) scanFinish(wsName model.WSName) *fsa.Action {
	return &fsa.Action{
		Type:    AssetScanFinishType,
		Payload: newWSPayload(wsName),
	}
}

type assetScanHandler struct {
	c                  <-chan *model.Asset
	queryCancel        context.CancelFunc
	assetUseCase       *usecase.Asset
	assetActionCreator *assetActionCreator
}

func (d *assetScanHandler) Do(action *fsa.Action, dispatch fsa.Dispatch) error {
	var payload assetScanRequestPayload
	if err := mapstructure.Decode(action.Payload, &payload); err != nil {
		return fmt.Errorf("failed to decode payload: %w", err)
	}

	b := context.Background()
	if d.c == nil || payload.Reset {
		if d.c != nil {
			d.queryCancel()
		}
		ctx, cancel := context.WithCancel(b)
		d.queryCancel = cancel
		c, err := d.assetUseCase.ListAsyncByQueries(ctx, payload.WorkSpaceName, payload.Queries)
		if err != nil {
			return err
		}
		d.c = c
	}

	var ret []*model.Asset
	for asset := range d.c {
		ret = append(ret, asset)
		if len(ret) >= payload.RequestNum {
			return dispatch(d.assetActionCreator.scanRunning(payload.WorkSpaceName, ret))
		}
	}
	if len(ret) > 0 {
		if err := dispatch(d.assetActionCreator.scanRunning(payload.WorkSpaceName, ret)); err != nil {
			return err
		}
	}
	return dispatch(d.assetActionCreator.scanFinish(payload.WorkSpaceName))
}

type assetHandlerCreator struct {
	assetUseCase       *usecase.Asset
	assetActionCreator *assetActionCreator
}

func newAssetHandlerCreator(
	assetUseCase *usecase.Asset,
) *assetHandlerCreator {
	return &assetHandlerCreator{
		assetUseCase:       assetUseCase,
		assetActionCreator: &assetActionCreator{},
	}
}

func (h *assetHandlerCreator) Scan() *assetScanHandler {
	return &assetScanHandler{
		assetUseCase:       h.assetUseCase,
		assetActionCreator: h.assetActionCreator,
	}
}
