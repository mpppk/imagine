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

type tagActionCreator struct{}

func (t *tagActionCreator) scan(wsName model.WSName, tags []*model.Tag) *fsa.Action {
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

func (t *tagActionCreator) save(wsName model.WSName, tags []*model.Tag) *fsa.Action {
	return &fsa.Action{
		Type: TagSaveType,
		Payload: &tagSavePayload{
			wsPayload: wsPayload{WorkSpaceName: wsName},
			Tags:      tags,
		},
	}
}

type tagScanHandler struct {
	tagUseCase *usecase.Tag
	action     *tagActionCreator
}

func (d *tagScanHandler) Do(action *fsa.Action, dispatch fsa.Dispatch) error {
	var payload tagRequestPayload
	if err := mapstructure.Decode(action.Payload, &payload); err != nil {
		return fmt.Errorf("failed to decode payload: %w", err)
	}

	tags, err := d.tagUseCase.List(payload.Name)
	if err != nil {
		return fmt.Errorf("failed to list tags: %w", err)
	}
	return dispatch(d.action.scan(payload.Name, tags))
}

type tagSaveHandler struct {
	tagUseCase *usecase.Tag
	action     *tagActionCreator
}

func (t *tagSaveHandler) Do(action *fsa.Action, dispatch fsa.Dispatch) error {
	var payload tagUpdatePayload
	if err := mapstructure.Decode(action.Payload, &payload); err != nil {
		return fmt.Errorf("failed to decode payload: %w", err)
	}
	fmt.Println(action)
	if err := t.tagUseCase.SetTags(payload.WorkSpaceName, payload.Tags); err != nil {
		return fmt.Errorf("failed to handle TagUpdate action: %w", err)
	}
	return dispatch(t.action.save(payload.WorkSpaceName, payload.Tags))
}

type tagHandlerCreator struct {
	tagUseCase *usecase.Tag
	action     *tagActionCreator
}

func newTagHandlerCreator(tagUseCase *usecase.Tag) *tagHandlerCreator {
	return &tagHandlerCreator{
		tagUseCase: tagUseCase,
		action:     &tagActionCreator{},
	}
}

func (h *tagHandlerCreator) Scan() *tagScanHandler {
	return &tagScanHandler{
		tagUseCase: h.tagUseCase,
		action:     h.action,
	}
}

func (h *tagHandlerCreator) Save() *tagSaveHandler {
	return &tagSaveHandler{
		tagUseCase: h.tagUseCase,
		action:     h.action,
	}
}
