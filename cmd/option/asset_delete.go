package option

import (
	"fmt"

	"github.com/mpppk/imagine/domain/model"
	"github.com/spf13/viper"
)

// AssetDeleteCmdConfig is config for eval command
type AssetDeleteCmdConfig struct {
	DB        string
	WorkSpace model.WSName
}

// NewAssetDeleteCmdConfigFromViper generate config for eval command from viper
func NewAssetDeleteCmdConfigFromViper(args []string) (*AssetDeleteCmdConfig, error) {
	var conf AssetDeleteCmdConfig
	if err := viper.Unmarshal(&conf); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config from viper: %w", err)
	}

	if err := conf.validate(); err != nil {
		return nil, fmt.Errorf("failed to create sum cmd config: %w", err)
	}

	return &conf, nil
}

func (c *AssetDeleteCmdConfig) validate() error {
	return nil
}
