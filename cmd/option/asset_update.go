package option

import (
	"fmt"

	"github.com/mpppk/imagine/domain/model"
	"github.com/spf13/viper"
)

// AssetAddCmdConfig is config for eval command
type AssetAddCmdConfig struct {
	DB        string
	WorkSpace model.WSName
}

// NewAssetAddCmdConfigFromViper generate config for eval command from viper
func NewAssetAddCmdConfigFromViper(args []string) (*AssetAddCmdConfig, error) {
	var conf AssetAddCmdConfig
	if err := viper.Unmarshal(&conf); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config from viper: %w", err)
	}

	if err := conf.validate(); err != nil {
		return nil, fmt.Errorf("failed to create sum cmd config: %w", err)
	}

	return &conf, nil
}

func (c *AssetAddCmdConfig) validate() error {
	return nil
}
