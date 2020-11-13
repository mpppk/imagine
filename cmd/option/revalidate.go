package option

import (
	"fmt"

	"github.com/mpppk/imagine/domain/model"

	"github.com/spf13/viper"
)

// RevalidateCmdConfig is config for eval command
type RevalidateCmdConfig struct {
	DB        string
	WorkSpace model.WSName
}

// NewRevalidateCmdConfigFromViper generate config for eval command from viper
func NewRevalidateCmdConfigFromViper(args []string) (*RevalidateCmdConfig, error) {
	var conf RevalidateCmdConfig
	if err := viper.Unmarshal(&conf); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config from viper: %w", err)
	}

	if err := conf.validate(); err != nil {
		return nil, fmt.Errorf("failed to create sum cmd config: %w", err)
	}

	return &conf, nil
}

func (c *RevalidateCmdConfig) validate() error {
	return nil
}
