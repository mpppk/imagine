package option

import (
	"fmt"

	"github.com/mpppk/imagine/domain/model"
	"github.com/spf13/viper"
)

// AssetUpdateCmdConfig is config for asset update command
type RawAssetUpdateCmdConfig struct {
	DB        string
	WorkSpace model.WSName
	New       bool
	Query     string
}

func (c *RawAssetUpdateCmdConfig) parse() (*AssetUpdateCmdConfig, error) {
	queries, err := parseQuery(c.Query)
	if err != nil {
		return nil, err
	}
	return &AssetUpdateCmdConfig{
		RawAssetUpdateCmdConfig: c,
		Queries:                 queries,
	}, nil
}

// AssetUpdateCmdConfig is config for eval command
type AssetUpdateCmdConfig struct {
	*RawAssetUpdateCmdConfig
	Queries []*model.Query
}

// NewAssetUpdateCmdConfigFromViper generate config for eval command from viper
func NewAssetUpdateCmdConfigFromViper() (*AssetUpdateCmdConfig, error) {
	var conf RawAssetUpdateCmdConfig
	if err := viper.Unmarshal(&conf); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config from viper: %w", err)
	}

	return conf.parse()
}
