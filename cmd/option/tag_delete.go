package option

import (
	"fmt"

	"github.com/mpppk/imagine/domain/model"

	"github.com/spf13/viper"
)

// TagDeleteCmdConfig is config for asset update command
type RawTagDeleteCmdConfig struct {
	DB        string
	WorkSpace model.WSName
	Query     string
}

func (c *RawTagDeleteCmdConfig) parse() (*TagDeleteCmdConfig, error) {
	queries, err := parseQuery(c.Query)
	if err != nil {
		return nil, err
	}
	return &TagDeleteCmdConfig{
		RawTagDeleteCmdConfig: c,
		Queries:               queries,
	}, nil
}

// TagDeleteCmdConfig is config for eval command
type TagDeleteCmdConfig struct {
	*RawTagDeleteCmdConfig
	Queries []*model.Query
}

// NewTagDeleteCmdConfigFromViper generate config for delete command
func NewTagDeleteCmdConfigFromViper(args []string) (*TagDeleteCmdConfig, error) {
	var conf RawTagDeleteCmdConfig
	if err := viper.Unmarshal(&conf); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config from viper: %w", err)
	}

	return conf.parse()
}
