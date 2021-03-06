package action

import (
	"fmt"

	"github.com/mpppk/imagine/domain/model"
	"github.com/mpppk/imagine/usecase"

	"github.com/mitchellh/mapstructure"
	fsa "github.com/mpppk/lorca-fsa"
)

const (
	boxPrefix                       = "BOUNDING_BOX/"
	BoxAssignRequestType   fsa.Type = boxPrefix + "ASSIGN/REQUEST"
	BoxAssignType          fsa.Type = boxPrefix + "ASSIGN"
	BoxUnAssignRequestType fsa.Type = boxPrefix + "UNASSIGN/REQUEST"
	BoxUnAssignType        fsa.Type = boxPrefix + "UNASSIGN"
	BoxModifyRequestType   fsa.Type = boxPrefix + "MODIFY/REQUEST"
	BoxModifyType          fsa.Type = boxPrefix + "MODIFY"
	BoxDeleteRequestType   fsa.Type = boxPrefix + "DELETE/REQUEST"
	BoxDeleteType          fsa.Type = boxPrefix + "DELETE"
)

type boxAssignRequestPayload struct {
	WsPayload `mapstructure:",squash"`
	Asset     *model.Asset       `json:"asset"`
	Box       *model.BoundingBox `json:"box"`
}

type boxUnAssignRequestPayload struct {
	WsPayload `mapstructure:",squash"`
	Asset     *model.Asset        `json:"asset"`
	BoxID     model.BoundingBoxID `json:"boxID"`
}

type boxModifyRequestPayload struct {
	WsPayload `mapstructure:",squash"`
	Asset     *model.Asset       `json:"asset"`
	Box       *model.BoundingBox `json:"box"`
}

type boxAssignPayload struct {
	WsPayload `mapstructure:",squash"`
	Asset     *model.Asset       `json:"asset"`
	Box       *model.BoundingBox `json:"box"`
}

type boxUnAssignPayload struct {
	WsPayload `mapstructure:",squash"`
	Asset     *model.Asset        `json:"asset"`
	BoxID     model.BoundingBoxID `json:"boxID"`
}

type boxDeleteRequestPayload struct {
	WsPayload `mapstructure:",squash"`
	AssetID   model.AssetID       `json:"assetID"`
	BoxID     model.BoundingBoxID `json:"boxID"`
}

type boxDeletePayload = boxDeleteRequestPayload

type boxActionCreator struct{}

func (a *boxActionCreator) assign(name model.WSName, asset *model.Asset, box *model.BoundingBox) *fsa.Action {
	return &fsa.Action{
		Type: BoxAssignType,
		Payload: &boxAssignPayload{
			WsPayload: WsPayload{WorkSpaceName: name},
			Asset:     asset,
			Box:       box,
		},
	}
}

func (a *boxActionCreator) unassign(name model.WSName, asset *model.Asset, boxID model.BoundingBoxID) *fsa.Action {
	return &fsa.Action{
		Type: BoxUnAssignType,
		Payload: &boxUnAssignPayload{
			WsPayload: WsPayload{WorkSpaceName: name},
			Asset:     asset,
			BoxID:     boxID,
		},
	}
}

func (a *boxActionCreator) modify(name model.WSName, asset *model.Asset, boxID model.BoundingBoxID) *fsa.Action {
	return &fsa.Action{
		Type: BoxModifyType,
		Payload: &boxUnAssignPayload{
			WsPayload: WsPayload{WorkSpaceName: name},
			Asset:     asset,
			BoxID:     boxID,
		},
	}
}

func (a *boxActionCreator) delete(name model.WSName, assetID model.AssetID, boxID model.BoundingBoxID) *fsa.Action {
	return &fsa.Action{
		Type: BoxDeleteType,
		Payload: &boxDeletePayload{
			WsPayload: WsPayload{WorkSpaceName: name},
			AssetID:   assetID,
			BoxID:     boxID,
		},
	}
}

type boxAssignRequestHandler struct {
	assetUseCase     usecase.Asset
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
	assetUseCase     usecase.Asset
	boxActionCreator *boxActionCreator
}

func (d *boxUnAssignRequestHandler) Do(action *fsa.Action, dispatch fsa.Dispatch) error {
	var payload boxUnAssignRequestPayload
	if err := mapstructure.Decode(action.Payload, &payload); err != nil {
		return fmt.Errorf("failed to decode payload: %w", err)
	}

	asset, err := d.assetUseCase.UnAssignBoundingBox(payload.WorkSpaceName, payload.Asset.ID, payload.BoxID)
	if err != nil {
		return fmt.Errorf("failed to unassgin bounding box from asset. boxID: %d, asset: %#v: %w", payload.BoxID, payload.Asset, err)
	}
	return dispatch(d.boxActionCreator.unassign(payload.WorkSpaceName, asset, payload.BoxID))
}

type boxModifyRequestHandler struct {
	assetUseCase     usecase.Asset
	boxActionCreator *boxActionCreator
}

func (d *boxModifyRequestHandler) Do(action *fsa.Action, dispatch fsa.Dispatch) error {
	var payload boxModifyRequestPayload
	if err := mapstructure.Decode(action.Payload, &payload); err != nil {
		return fmt.Errorf("failed to decode payload: %w", err)
	}

	asset, err := d.assetUseCase.ModifyBoundingBox(payload.WorkSpaceName, payload.Asset.ID, payload.Box)
	if err != nil {
		return fmt.Errorf("failed to modify bounding box from asset. boxID: %d, asset: %#v: %w", payload.Box.ID, payload.Asset, err)
	}
	return dispatch(d.boxActionCreator.modify(payload.WorkSpaceName, asset, payload.Box.ID))
}

type boxDeleteRequestHandler struct {
	assetUseCase     usecase.Asset
	boxActionCreator *boxActionCreator
}

func (d *boxDeleteRequestHandler) Do(action *fsa.Action, dispatch fsa.Dispatch) error {
	var payload boxDeleteRequestPayload
	if err := mapstructure.Decode(action.Payload, &payload); err != nil {
		return fmt.Errorf("failed to decode payload: %w", err)
	}

	if err := d.assetUseCase.DeleteBoundingBox(payload.WorkSpaceName, payload.AssetID, payload.BoxID); err != nil {
		return fmt.Errorf("failed to modify bounding box from asset. boxID: %d, asset: %#v: %w", payload.BoxID, payload.AssetID, err)
	}
	return dispatch(d.boxActionCreator.delete(payload.WorkSpaceName, payload.AssetID, payload.BoxID))
}

type boxHandlerCreator struct {
	assetUseCase     usecase.Asset
	boxActionCreator *boxActionCreator
}

func newBoxHandlerCreator(
	assetUseCase usecase.Asset,
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

func (h *boxHandlerCreator) Modify() *boxModifyRequestHandler {
	return &boxModifyRequestHandler{
		assetUseCase:     h.assetUseCase,
		boxActionCreator: h.boxActionCreator,
	}
}

func (h *boxHandlerCreator) Delete() *boxDeleteRequestHandler {
	return &boxDeleteRequestHandler{
		assetUseCase:     h.assetUseCase,
		boxActionCreator: h.boxActionCreator,
	}
}
