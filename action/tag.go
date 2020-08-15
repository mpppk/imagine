package action

import (
	"fmt"

	"github.com/mpppk/imagine/domain/model"

	"github.com/mitchellh/mapstructure"
	"github.com/mpppk/imagine/usecase"
	fsa "github.com/mpppk/lorca-fsa/lorca-fsa"
)

type TagRequestHandler struct {
	tagUseCase *usecase.Tag
}

func NewTagRequestHandler(assetUseCase *usecase.Tag) *TagRequestHandler {
	return &TagRequestHandler{tagUseCase: assetUseCase}
}

type TagRequestPayload struct {
	WSPayload `mapstructure:",squash"`
}

type TagScanPayload struct {
	*WSPayload
	Tags []*model.Tag `json:"assets"`
}

func newTagScanAction(wsName model.WSName, tags []*model.Tag) *fsa.Action {
	return &fsa.Action{
		Type: ServerTagScanType,
		Payload: &TagScanPayload{
			WSPayload: newWSPayload(wsName),
			Tags:      tags,
		},
	}
}

func (d *TagRequestHandler) Do(action *fsa.Action, dispatch fsa.Dispatch) error {
	var payload TagRequestPayload
	if err := mapstructure.Decode(action.Payload, &payload); err != nil {
		return fmt.Errorf("failed to decode payload: %w", err)
	}

	tags, err := d.tagUseCase.List(payload.WorkSpaceName)
	if err != nil {
		return fmt.Errorf("failed to list tags: %w", err)
	}
	return dispatch(newTagScanAction(payload.WorkSpaceName, tags))
}
