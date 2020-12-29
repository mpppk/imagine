package repository

import (
	"fmt"

	"github.com/mpppk/imagine/domain/model"
)

type Client struct {
	Asset     Asset
	Tag       Tag
	WorkSpace WorkSpace
	Meta      Meta
}

func NewClient(asset Asset, tag Tag, workspace WorkSpace, meta Meta) *Client {
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
		return fmt.Errorf("failed to initialize asset repository: %w", err)
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
