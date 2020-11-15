package option

import (
	"fmt"

	"github.com/mpppk/imagine/domain/model"

	"github.com/spf13/viper"
)

// MigrateCmdConfig is config for eval command
type MigrateCmdConfig struct {
	DB        string
	WorkSpace model.WSName
}

// NewMigrateCmdConfigFromViper generate config for eval command from viper
func NewMigrateCmdConfigFromViper(args []string) (*MigrateCmdConfig, error) {
	var conf MigrateCmdConfig
	if err := viper.Unmarshal(&conf); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config from viper: %w", err)
	}

	if err := conf.validate(); err != nil {
		return nil, fmt.Errorf("failed to create sum cmd config: %w", err)
	}

	return &conf, nil
}

func (c *MigrateCmdConfig) validate() error {
	return nil
}
