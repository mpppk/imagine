package cmd

import (
	"fmt"
	"log"
	"net/http"
	"os"

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
func NewRootCmd(fs afero.Fs) (*cobra.Command, error) {
	pPreRunE := func(cmd *cobra.Command, args []string) error {
		conf, err := option.NewRootCmdConfigFromViper()
		if err != nil {
			return err
		}
		util.InitializeLog(conf.Verbose)
		return nil
	}

	cmd := &cobra.Command{
		Use:               "imagine",
		Short:             "imagine",
		SilenceErrors:     true,
		SilenceUsage:      true,
		PersistentPreRunE: pPreRunE,
		RunE: func(cmd *cobra.Command, args []string) error {
			//devMode := false
			//if len(os.Args) > 1 && os.Args[1] == "dev" {
			//	devMode = true
			//}
			devMode := true

			db, err := bbolt.Open("test.db", 0600, nil)
			if err != nil {
				return fmt.Errorf("failed to open DB: %w", err)
			}

			logger := colog.NewCoLog(os.Stdout, "", 0).NewLogger()

			handlers := registry.NewHandlers(db)
			config := &fsa.LorcaConfig{
				AppName:          "imagine",
				Url:              "http://localhost:3000",
				Width:            1080,
				Height:           720,
				EnableExtensions: devMode,
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

			http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("/"))))

			go func() {
				if err := http.ListenAndServe(":1323", nil); err != nil {
					log.Fatal("ListenAndServe: ", err)
				}
			}()

			fsa.Wait(ui)
			return nil
		},
	}

	if err := registerSubCommands(fs, cmd); err != nil {
		return nil, err
	}

	if err := registerFlags(cmd); err != nil {
		return nil, err
	}

	return cmd, nil
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
				Name:         "config",
				IsPersistent: true,
				Usage:        "config file (default is $HOME/.imagine.yaml)",
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
	rootCmd, err := NewRootCmd(afero.NewOsFs())
	if err != nil {
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
