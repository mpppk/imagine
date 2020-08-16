package action

import (
	"fmt"

	"github.com/mpppk/imagine/domain/model"

	"github.com/mitchellh/mapstructure"
	"github.com/mpppk/imagine/usecase"
	fsa "github.com/mpppk/lorca-fsa/lorca-fsa"
)

const (
	tagPrefix               = "TAG/"
	TagRequestType fsa.Type = tagPrefix + "REQUEST"
	TagScanType    fsa.Type = tagPrefix + "SCAN"
	TagSaveType    fsa.Type = tagPrefix + "SAVE"
)

type tagRequestHandler struct {
	tagUseCase *usecase.Tag
}

type tagRequestPayload = model.WorkSpace

type tagScanPayload struct {
	wsPayload `mapstructure:",squash"`
	Tags      []*model.Tag `json:"tags"`
}

type tagSavePayload struct {
	wsPayload `mapstructure:",squash"`
	Tags      []*model.Tag `json:"tags"`
}

type tagUpdatePayload struct {
	wsPayload `mapstructure:",squash"`
	Tags      []*model.Tag `json:"tags"`
}

func newTagScanAction(wsName model.WSName, tags []*model.Tag) *fsa.Action {
	if tags == nil {
		tags = []*model.Tag{}
	}
	return &fsa.Action{
		Type: TagScanType,
		Payload: &tagScanPayload{
			wsPayload: wsPayload{WorkSpaceName: wsName},
			Tags:      tags,
		},
	}
}

func newTagSaveAction(wsName model.WSName, tags []*model.Tag) *fsa.Action {
	return &fsa.Action{
		Type: TagSaveType,
		Payload: &tagSavePayload{
			wsPayload: wsPayload{WorkSpaceName: wsName},
			Tags:      tags,
		},
	}
}

func (d *tagRequestHandler) Do(action *fsa.Action, dispatch fsa.Dispatch) error {
	var payload tagRequestPayload
	if err := mapstructure.Decode(action.Payload, &payload); err != nil {
		return fmt.Errorf("failed to decode payload: %w", err)
	}

	tags, err := d.tagUseCase.List(payload.Name)
	if err != nil {
		return fmt.Errorf("failed to list tags: %w", err)
	}
	return dispatch(newTagScanAction(payload.Name, tags))
}

type tagUpdateHandler struct {
	tagUseCase *usecase.Tag
}

func (d *tagUpdateHandler) Do(action *fsa.Action, dispatch fsa.Dispatch) error {
	var payload tagUpdatePayload
	if err := mapstructure.Decode(action.Payload, &payload); err != nil {
		return fmt.Errorf("failed to decode payload: %w", err)
	}
	fmt.Println(action)
	if err := d.tagUseCase.SetTags(payload.WorkSpaceName, payload.Tags); err != nil {
		return fmt.Errorf("failed to handle TagUpdate action: %w", err)
	}
	return dispatch(newTagSaveAction(payload.WorkSpaceName, payload.Tags))
}
