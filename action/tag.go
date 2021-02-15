package action

import (
	"fmt"

	"github.com/mpppk/imagine/domain/service/assetsvc/tagsvc"

	"github.com/mpppk/imagine/usecase"

	"github.com/mpppk/imagine/domain/model"

	"github.com/mitchellh/mapstructure"
	fsa "github.com/mpppk/lorca-fsa"
)

const (
	tagPrefix                  = "TAG/"
	TagScanResultType fsa.Type = tagPrefix + "SCAN/RESULT"
	TagSaveType       fsa.Type = tagPrefix + "SAVE"
	TagUpdateType     fsa.Type = tagPrefix + "UPDATE"
)

type tagRequestPayload = model.WorkSpace

type tagScanPayload struct {
	WsPayload `mapstructure:",squash"`
	Tags      []*model.Tag `json:"tags"`
}

type tagSavePayload struct {
	WsPayload `mapstructure:",squash"`
	Tags      []*model.Tag `json:"tags"`
}

type tagUpdatePayload struct {
	WsPayload `mapstructure:",squash"`
	Tags      []*model.Tag `json:"tags"`
}

type tagActionCreator struct{}

func (t *tagActionCreator) scan(wsName model.WSName, tags []*model.Tag) *fsa.Action {
	if tags == nil {
		tags = []*model.Tag{}
	}
	return &fsa.Action{
		Type: TagScanResultType,
		Payload: &tagScanPayload{
			WsPayload: WsPayload{WorkSpaceName: wsName},
			Tags:      tags,
		},
	}
}

func (t *tagActionCreator) save(wsName model.WSName, tags []*model.Tag) *fsa.Action {
	return &fsa.Action{
		Type: TagSaveType,
		Payload: &tagSavePayload{
			WsPayload: WsPayload{WorkSpaceName: wsName},
			Tags:      tags,
		},
	}
}

type tagScanHandler struct {
	tagUseCase usecase.Tag
	action     *tagActionCreator
}

func (d *tagScanHandler) Do(action *fsa.Action, dispatch fsa.Dispatch) error {
	var payload tagRequestPayload
	if err := mapstructure.Decode(action.Payload, &payload); err != nil {
		return fmt.Errorf("failed to decode payload: %w", err)
	}

	if !payload.Name.IsValid() {
		return fmt.Errorf("invalid workspace name: %q", payload.Name)
	}

	tags, err := d.tagUseCase.List(payload.Name)
	if err != nil {
		return fmt.Errorf("failed to list tags: %w", err)
	}
	return dispatch(d.action.scan(payload.Name, tags))
}

type tagSaveHandler struct {
	tagUseCase usecase.Tag
	action     *tagActionCreator
}

func (t *tagSaveHandler) Do(action *fsa.Action, dispatch fsa.Dispatch) error {
	var payload tagUpdatePayload
	if err := mapstructure.Decode(action.Payload, &payload); err != nil {
		return fmt.Errorf("failed to decode payload: %w", err)
	}

	tagNames := tagsvc.ToTagNames(payload.Tags)
	if _, err := t.tagUseCase.SetTags(payload.WorkSpaceName, tagNames); err != nil {
		return fmt.Errorf("failed to handle TagUpdate action: %w", err)
	}
	return dispatch(t.action.save(payload.WorkSpaceName, payload.Tags))
}

type tagHandlerCreator struct {
	tagUseCase usecase.Tag
	action     *tagActionCreator
}

func newTagHandlerCreator(tagUseCase usecase.Tag) *tagHandlerCreator {
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
