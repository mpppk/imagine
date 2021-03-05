package option

import (
	"fmt"

	"github.com/mpppk/imagine/domain/model"
	"github.com/spf13/viper"
)

// AssetAddCmdConfig is config for asset update command
type RawAssetAddCmdConfig struct {
	DB        string
	WorkSpace model.WSName
	New       bool
	Query     string
}

func (c *RawAssetAddCmdConfig) parse() (*AssetAddCmdConfig, error) {
	queries, err := parseQuery(c.Query)
	if err != nil {
		return nil, err
	}
	return &AssetAddCmdConfig{
		RawAssetAddCmdConfig: c,
		Queries:              queries,
	}, nil
}

// AssetAddCmdConfig is config for eval command
type AssetAddCmdConfig struct {
	*RawAssetAddCmdConfig
	Queries []*model.Query
}

// NewAssetAddCmdConfigFromViper generate config for eval command from viper
func NewAssetAddCmdConfigFromViper(args []string) (*AssetAddCmdConfig, error) {
	var conf RawAssetAddCmdConfig
	if err := viper.Unmarshal(&conf); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config from viper: %w", err)
	}

	return conf.parse()
}
