package action

import (
	"fmt"

	"github.com/mpppk/imagine/domain/model"

	"github.com/mitchellh/mapstructure"
	"github.com/mpppk/imagine/usecase"
	fsa "github.com/mpppk/lorca-fsa/lorca-fsa"
)

const (
	boxPrefix                       = "BOUNDING_BOX/"
	BoxAssignRequestType   fsa.Type = boxPrefix + "ASSIGN/REQUEST"
	BoxAssignType          fsa.Type = boxPrefix + "ASSIGN"
	BoxUnAssignRequestType fsa.Type = boxPrefix + "UNASSIGN/REQUEST"
	BoxUnAssignType        fsa.Type = boxPrefix + "UNASSIGN"
)

type boxAssignRequestPayload struct {
	wsPayload `mapstructure:",squash"`
	Asset     *model.Asset       `json:"asset"`
	Box       *model.BoundingBox `json:"box"`
}

type boxUnAssignRequestPayload struct {
	wsPayload `mapstructure:",squash"`
	Asset     *model.Asset        `json:"asset"`
	BoxID     model.BoundingBoxID `json:"boxID"`
}

type boxAssignPayload struct {
	wsPayload `mapstructure:",squash"`
	Asset     *model.Asset       `json:"asset"`
	Box       *model.BoundingBox `json:"box"`
}

type boxUnAssignPayload struct {
	wsPayload `mapstructure:",squash"`
	Asset     *model.Asset        `json:"asset"`
	BoxID     model.BoundingBoxID `json:"boxID"`
}

type boxActionCreator struct{}

func (a *boxActionCreator) assign(name model.WSName, asset *model.Asset, box *model.BoundingBox) *fsa.Action {
	return &fsa.Action{
		Type: BoxAssignType,
		Payload: &boxAssignPayload{
			wsPayload: wsPayload{WorkSpaceName: name},
			Asset:     asset,
			Box:       box,
		},
	}
}

func (a *boxActionCreator) unassign(name model.WSName, asset *model.Asset, boxID model.BoundingBoxID) *fsa.Action {
	return &fsa.Action{
		Type: BoxUnAssignType,
		Payload: &boxUnAssignPayload{
			wsPayload: wsPayload{WorkSpaceName: name},
			Asset:     asset,
			BoxID:     boxID,
		},
	}
}

type boxAssignRequestHandler struct {
	c                <-chan *model.Asset
	assetUseCase     *usecase.Asset
	boxActionCreator *boxActionCreator
}

func (d *boxAssignRequestHandler) Do(action *fsa.Action, dispatch fsa.Dispatch) error {
	var payload boxAssignRequestPayload
	if err := mapstructure.Decode(action.Payload, &payload); err != nil {
		return fmt.Errorf("failed to decode payload: %w", err)
	}
	asset, err := d.assetUseCase.AssignBoundingBox(payload.WorkSpaceName, payload.Asset.ID, payload.Box)
	if err != nil {
		return fmt.Errorf("failed to assgin bounding box to asset. box: %#v, asset: %#v: %w", payload.Box, payload.Asset, err)
	}
	return dispatch(d.boxActionCreator.assign(payload.WorkSpaceName, asset, payload.Box))
}

type boxUnAssignRequestHandler struct {
	c                <-chan *model.Asset
	assetUseCase     *usecase.Asset
	boxActionCreator *boxActionCreator
}

func (d *boxUnAssignRequestHandler) Do(action *fsa.Action, dispatch fsa.Dispatch) error {
	var payload boxUnAssignRequestPayload
	if err := mapstructure.Decode(action.Payload, &payload); err != nil {
		return fmt.Errorf("failed to decode payload: %w", err)
	}

	fmt.Printf("unassign payload: %#v\n", payload)

	asset, err := d.assetUseCase.UnAssignBoundingBox(payload.WorkSpaceName, payload.Asset.ID, payload.BoxID)
	if err != nil {
		return fmt.Errorf("failed to unassgin bounding box from asset. boxID: %d, asset: %#v: %w", payload.BoxID, payload.Asset, err)
	}
	return dispatch(d.boxActionCreator.unassign(payload.WorkSpaceName, asset, payload.BoxID))
}

type boxHandlerCreator struct {
	assetUseCase     *usecase.Asset
	boxActionCreator *boxActionCreator
}

func newBoxHandlerCreator(
	assetUseCase *usecase.Asset,
) *boxHandlerCreator {
	return &boxHandlerCreator{
		assetUseCase:     assetUseCase,
		boxActionCreator: &boxActionCreator{},
	}
}

func (h *boxHandlerCreator) Assign() *boxAssignRequestHandler {
	return &boxAssignRequestHandler{
		assetUseCase:     h.assetUseCase,
		boxActionCreator: h.boxActionCreator,
	}
}

func (h *boxHandlerCreator) UnAssign() *boxUnAssignRequestHandler {
	return &boxUnAssignRequestHandler{
		assetUseCase:     h.assetUseCase,
		boxActionCreator: h.boxActionCreator,
	}
}
