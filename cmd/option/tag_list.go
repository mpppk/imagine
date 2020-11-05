package option

import (
	"fmt"

	"github.com/mpppk/imagine/domain/model"

	"github.com/spf13/viper"
)

// TagListCmdConfig is config for eval command
type TagListCmdConfig struct {
	DB        string
	WorkSpace model.WSName
}

// NewTagListCmdConfigFromViper generate config for eval command from viper
func NewTagListCmdConfigFromViper(args []string) (*TagListCmdConfig, error) {
	var conf TagListCmdConfig
	if err := viper.Unmarshal(&conf); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config from viper: %w", err)
	}

	if err := conf.validate(); err != nil {
		return nil, fmt.Errorf("failed to create sum cmd config: %w", err)
	}

	return &conf, nil
}

func (c *TagListCmdConfig) validate() error {
	return nil
}
