package option

import (
	"fmt"

	"github.com/mpppk/imagine/domain/model"

	"github.com/spf13/viper"
)

// ExportCmdConfig is config for eval command
type ExportCmdConfig struct {
	DB        string
	WorkSpace model.WSName
}

// NewExportCmdConfigFromViper generate config for eval command from viper
func NewExportCmdConfigFromViper(args []string) (*ExportCmdConfig, error) {
	var conf ExportCmdConfig
	if err := viper.Unmarshal(&conf); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config from viper: %w", err)
	}

	if err := conf.validate(); err != nil {
		return nil, fmt.Errorf("failed to create sum cmd config: %w", err)
	}

	return &conf, nil
}

func (c *ExportCmdConfig) validate() error {
	return nil
}
