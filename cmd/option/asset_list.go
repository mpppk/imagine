package option

import (
	"fmt"

	"github.com/mpppk/imagine/domain/model"

	"github.com/spf13/viper"
)

// AssetListCmdConfig is config for eval command
type AssetListCmdConfig struct {
	DB        string
	WorkSpace model.WSName
	Format    string
}

// NewAssetListCmdConfigFromViper generate config for eval command from viper
func NewAssetListCmdConfigFromViper(args []string) (*AssetListCmdConfig, error) {
	var conf AssetListCmdConfig
	if err := viper.Unmarshal(&conf); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config from viper: %w", err)
	}

	if err := conf.validate(); err != nil {
		return nil, fmt.Errorf("failed to create sum cmd config: %w", err)
	}

	return &conf, nil
}

func (c *AssetListCmdConfig) validate() error {
	return nil
}
