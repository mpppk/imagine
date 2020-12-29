package cmd

import (
	"fmt"
	"os"

	"github.com/mpppk/imagine/infra"

	"github.com/mpppk/imagine/registry"
	"github.com/mpppk/imagine/util"

	"github.com/mpppk/imagine/cmd/option"

	"github.com/spf13/afero"

	"github.com/mitchellh/go-homedir"
	_ "github.com/mpppk/imagine/statik"
	fsa "github.com/mpppk/lorca-fsa/lorca-fsa"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var RootCmd, rootCmdErr = NewRootCmd()

func NewRootCmd() (*cobra.Command, error) {
	var rootCmd = &cobra.Command{
		Use:           "imagine",
		Short:         "imagine",
		SilenceErrors: true,
		SilenceUsage:  true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			cmd.SetOut(os.Stdout)
			conf, err := option.NewRootCmdConfigFromViper()
			if err != nil {
				return err
			}
			util.InitializeLog(conf.Verbose)

			usecases, err := registry.NewBoltUseCasesWithDBPath(conf.DB)
			if err != nil {
				return fmt.Errorf("failed to create usecases instance: %w", err)
			}
			defer func() {
				if err := usecases.Close(); err != nil {
					panic(err)
				}
			}()
			if err := usecases.Client.Init(); err != nil {
				return fmt.Errorf("failed to initialize meta repository: %w", err)
			}

			// for debug. set version for test
			//v := semver.MustParse("0.0.1")
			//if err := client.Meta.SetDBVersion(&v); err != nil {
			//	return err
			//}

			if err := usecases.Migration.Migrate(); err != nil {
				return err
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			conf, err := option.NewRootCmdConfigFromViper()
			if err != nil {
				return err
			}

			logger := util.GetLogger()

			handlers, closer, err := registry.NewHandlersWithDBPath(conf.DB)
			if err != nil {
				return err
			}
			defer func() {
				if err := closer.Close(); err != nil {
					panic(err)
				}
			}()

			if !conf.Dev {
				s, err := infra.NewHtmlServer(conf.UiPort)
				if err != nil {
					return err
				}

				go func() {
					if err := s.ListenAndServe(); err != nil {
						panic(err)
					}
				}()
			}

			config := &fsa.LorcaConfig{
				AppName:          "imagine",
				Url:              fmt.Sprintf("localhost:%d", conf.UiPort),
				Width:            1080,
				Height:           720,
				EnableExtensions: conf.Dev,
				Handlers:         handlers,
				Logger:           logger,
			}

			ui, err := fsa.Start(config)
			if err != nil {
				panic(err)
			}
			defer func() {
				if err := ui.Close(); err != nil {
					panic(err)
				}
			}()

			fsa.Wait(ui)
			return nil
		},
	}
	if err := registerFlags(rootCmd); err != nil {
		return nil, fmt.Errorf("failed to register flags")
	}
	fs := afero.NewOsFs()
	if err := registerSubCommands(fs, rootCmd); err != nil {
		panic(err)
	}

	return rootCmd, nil
}

func registerSubCommands(fs afero.Fs, cmd *cobra.Command) error {
	var subCmds []*cobra.Command
	for _, cmdGen := range cmdGenerators {
		subCmd, err := cmdGen(fs)
		if err != nil {
			return err
		}
		subCmds = append(subCmds, subCmd)
	}
	cmd.AddCommand(subCmds...)
	return nil
}

func registerFlags(cmd *cobra.Command) error {
	flags := []option.Flag{
		&option.StringFlag{
			BaseFlag: &option.BaseFlag{
				Name:         "db",
				Usage:        "db file path",
				IsRequired:   true,
				IsPersistent: true,
			},
		},
		&option.StringFlag{
			BaseFlag: &option.BaseFlag{
				Name:         "config",
				IsPersistent: true,
				Usage:        "config file (default is $HOME/.imagine.yaml)",
			}},
		&option.BoolFlag{
			BaseFlag: &option.BaseFlag{
				Name:  "dev",
				Usage: "Launch as developer mode",
			}},
		&option.BoolFlag{
			BaseFlag: &option.BaseFlag{
				Name:         "verbose",
				Shorthand:    "v",
				IsPersistent: true,
				Usage:        "Show more logs",
			}},
		&option.UintFlag{
			BaseFlag: &option.BaseFlag{
				Name:         "ui-port",
				IsPersistent: true,
				Usage:        "port of ui server",
			},
			Value: 3001,
		},
		&option.UintFlag{
			BaseFlag: &option.BaseFlag{
				Name:         "asset-port",
				IsPersistent: true,
				Usage:        "port of asset server (WIP: currently does not work)",
			},
			Value: 1323,
		},
	}
	return option.RegisterFlags(cmd, flags)
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the RootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Print(util.PrettyPrintError(err))
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	if rootCmdErr != nil {
		panic(rootCmdErr)
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".imagine" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".imagine")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
