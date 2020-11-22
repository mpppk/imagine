package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/mpppk/imagine/usecase"

	"github.com/blang/semver/v4"

	"github.com/comail/colog"

	"go.etcd.io/bbolt"

	"github.com/mpppk/imagine/registry"
	"github.com/mpppk/imagine/util"

	"github.com/mpppk/imagine/cmd/option"

	"github.com/spf13/afero"

	"github.com/mitchellh/go-homedir"
	fsa "github.com/mpppk/lorca-fsa/lorca-fsa"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// NewRootCmd generate root cmd
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

		db, err := bbolt.Open(conf.DB, 0600, nil)
		if err != nil {
			return fmt.Errorf("failed to open DB: %w", err)
		}
		defer func() {
			if err := db.Close(); err != nil {
				panic(err)
			}
		}()

		client := registry.NewBoltClient(db)
		if err := client.Meta.Init(); err != nil {
			return fmt.Errorf("failed to initialize meta repository: %w", err)
		}

		// for debug. set version for test
		//v := semver.MustParse("0.0.1")
		//if err := client.Meta.SetDBVersion(&v); err != nil {
		//	return err
		//}

		dbV, ok, err := client.Meta.GetDBVersion()
		if err != nil {
			return fmt.Errorf("failed to get db version: %w", err)
		}

		appV := semver.MustParse(util.Version)
		if !ok {
			if err := client.Meta.SetDBVersion(&appV); err != nil {
				return err
			}
			log.Printf("info: versions: db:%s app:%s", "empty→"+appV.String(), appV.String())
		} else {
			log.Printf("info: versions: db:%s app:%s", dbV.String(), appV.String())
		}

		migrationUseCase := usecase.NewMigration(client)
		if err := migrationUseCase.Migrate(dbV); err != nil {
			return err
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		conf, err := option.NewRootCmdConfigFromViper()
		if err != nil {
			return err
		}

		db, err := bbolt.Open(conf.DB, 0600, nil)
		if err != nil {
			return fmt.Errorf("failed to open DB: %w", err)
		}

		defer func() {
			if err := db.Close(); err != nil {
				panic(err)
			}
		}()

		logger := colog.NewCoLog(os.Stdout, "", 0).NewLogger()

		handlers := registry.NewHandlers(db)
		config := &fsa.LorcaConfig{
			AppName:          "imagine",
			Url:              conf.UiURL,
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
		&option.StringFlag{
			BaseFlag: &option.BaseFlag{
				Name:  "ui-url",
				Usage: "URL of front end server",
			},
			Value: "localhost:3000", // FIXME embedded
		},
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
	}
	return option.RegisterFlags(cmd, flags)
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	fs := afero.NewOsFs()
	if err := registerSubCommands(fs, rootCmd); err != nil {
		panic(err)
	}

	if err := registerFlags(rootCmd); err != nil {
		panic(err)
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Print(util.PrettyPrintError(err))
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
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
