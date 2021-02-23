package client

import (
	"fmt"

	"github.com/mpppk/imagine/domain/query"

	"github.com/mpppk/imagine/domain/repository"

	"github.com/mpppk/imagine/domain/model"
)

type TagRepository = repository.Tag
type TagQuery = query.Tag

type Tag struct {
	TagRepository
	TagQuery
}

func NewTag(r repository.Tag, q query.Tag) *Tag {
	return &Tag{r, q}
}

func (t *Tag) Init(ws model.WSName) error {
	if err := t.TagRepository.Init(ws); err != nil {
		return err
	}
	if err := t.TagQuery.Init(ws); err != nil {
		return err
	}
	return nil
}

type Client struct {
	Asset     repository.Asset
	Tag       *Tag
	WorkSpace repository.WorkSpace
	Meta      repository.Meta
}

func New(asset repository.Asset, tag *Tag, workspace repository.WorkSpace, meta repository.Meta) *Client {
	return &Client{
		Asset:     asset,
		Tag:       tag,
		WorkSpace: workspace,
		Meta:      meta,
	}
}

func (c *Client) Init() error {
	if err := c.WorkSpace.Init(); err != nil {
		return fmt.Errorf("failed to initialize asset repository: %w", err)
	}
	if err := c.Meta.Init(); err != nil {
		return fmt.Errorf("failed to initialize asset repository: %w", err)
	}
	return nil
}

func (c *Client) initWorkSpace(ws model.WSName) error {
	if err := c.Asset.Init(ws); err != nil {
		return fmt.Errorf("failed to initialize asset repository: %w", err)
	}
	if err := c.Tag.Init(ws); err != nil {
		return fmt.Errorf("failed to initialize tag repository: %w", err)
	}
	return nil
}

func (c *Client) CreateWorkSpace(ws model.WSName) (*model.WorkSpace, error) {
	workspace := &model.WorkSpace{Name: ws}
	if err := c.WorkSpace.Add(workspace); err != nil {
		return nil, fmt.Errorf("failed to create default workspace: %w", err)
	}
	if err := c.initWorkSpace(ws); err != nil {
		return nil, fmt.Errorf("failed to initialize workspace: %w", err)
	}
	return workspace, nil
}

func (c *Client) Close() error {
	return c.Asset.Close()
}
