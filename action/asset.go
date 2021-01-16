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
	Count  int            `json:"count"`
}

type assetScanResultPayload struct {
	*wsPayload
	Count int `json:"count"`
}

type assetActionCreator struct{}

func (a *assetActionCreator) scanRunning(wsName model.WSName, assets []*model.Asset, count int) *fsa.Action {
	return &fsa.Action{
		Type: AssetScanRunningType,
		Payload: &assetScanRunningPayload{
			wsPayload: newWSPayload(wsName),
			Assets:    assets,
			Count:     count,
		},
	}
}

func (a *assetActionCreator) scanFinish(wsName model.WSName, count int) *fsa.Action {
	return &fsa.Action{
		Type: AssetScanFinishType,
		Payload: &assetScanResultPayload{
			wsPayload: newWSPayload(wsName),
			Count:     count,
		},
	}
}

type assetScanHandler struct {
	c                  <-chan *model.Asset
	cnt                int
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
			d.cnt += len(ret)
			return dispatch(d.assetActionCreator.scanRunning(payload.WorkSpaceName, ret, d.cnt))
		}
	}

	if len(ret) > 0 {
		d.cnt += len(ret)
		return dispatch(d.assetActionCreator.scanRunning(payload.WorkSpaceName, ret, d.cnt))
	} else {
		return dispatch(d.assetActionCreator.scanFinish(payload.WorkSpaceName, d.cnt))
	}
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
