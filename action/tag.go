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

type TagRequestHandler struct {
	tagUseCase *usecase.Tag
}

func NewTagRequestHandler(tagUseCase *usecase.Tag) *TagRequestHandler {
	return &TagRequestHandler{tagUseCase: tagUseCase}
}

type TagRequestPayload = model.WorkSpace

type TagScanPayload struct {
	wsPayload `mapstructure:",squash"`
	Tags      []*model.Tag `json:"tags"`
}

type TagSavePayload struct {
	wsPayload `mapstructure:",squash"`
	Tags      []*model.Tag `json:"tags"`
}

type TagUpdatePayload struct {
	wsPayload `mapstructure:",squash"`
	Tags      []*model.Tag `json:"tags"`
}

func newTagScanAction(wsName model.WSName, tags []*model.Tag) *fsa.Action {
	if tags == nil {
		tags = []*model.Tag{}
	}
	return &fsa.Action{
		Type: TagScanType,
		Payload: &TagScanPayload{
			wsPayload: wsPayload{WorkSpaceName: wsName},
			Tags:      tags,
		},
	}
}

func newTagSaveAction(wsName model.WSName, tags []*model.Tag) *fsa.Action {
	return &fsa.Action{
		Type: TagSaveType,
		Payload: &TagSavePayload{
			wsPayload: wsPayload{WorkSpaceName: wsName},
			Tags:      tags,
		},
	}
}

func (d *TagRequestHandler) Do(action *fsa.Action, dispatch fsa.Dispatch) error {
	var payload TagRequestPayload
	if err := mapstructure.Decode(action.Payload, &payload); err != nil {
		return fmt.Errorf("failed to decode payload: %w", err)
	}

	tags, err := d.tagUseCase.List(payload.Name)
	if err != nil {
		return fmt.Errorf("failed to list tags: %w", err)
	}
	return dispatch(newTagScanAction(payload.Name, tags))
}

type TagUpdateHandler struct {
	tagUseCase *usecase.Tag
}

func (d *TagUpdateHandler) Do(action *fsa.Action, dispatch fsa.Dispatch) error {
	var payload TagUpdatePayload
	if err := mapstructure.Decode(action.Payload, &payload); err != nil {
		return fmt.Errorf("failed to decode payload: %w", err)
	}
	fmt.Println(action)
	if err := d.tagUseCase.SetTags(payload.WorkSpaceName, payload.Tags); err != nil {
		return fmt.Errorf("failed to handle TagUpdate action: %w", err)
	}
	return dispatch(newTagSaveAction(payload.WorkSpaceName, payload.Tags))
}

func NewTagUpdateHandler(tagUseCase *usecase.Tag) *TagUpdateHandler {
	return &TagUpdateHandler{tagUseCase: tagUseCase}
}
