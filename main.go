package main

import (
	"os"

	"gitea.antoine-langlois.net/datahearth/config-mapper/internal/actions"
	cmConfig "gitea.antoine-langlois.net/datahearth/config-mapper/internal/config"
	cmLog "gitea.antoine-langlois.net/datahearth/config-mapper/internal/logger"
	"github.com/urfave/cli/v2"
)

var app = &cli.App{
	Version:     "v1.0.0-alpha",
	Name:        "config-mapper",
	Description: "Manage your systems configurations",
	UsageText: `config-mapper aims to help you manage your configurations between systems
with a single configuration file.`,
	Authors: []*cli.Author{
		{
			Name:  "Antoine Langlois",
			Email: "antoine.l@antoine-langlois.net",
		},
	},
	Suggest:              true,
	EnableBashCompletion: true,
	Flags:                []cli.Flag{},
	Commands: []*cli.Command{
		{
			Name:        "setup",
			Description: "Setup the configuration file",
			Action: func(ctx *cli.Context) error {
				if err := actions.Setup(); err != nil {
					return err
				}

				logger.Info("Configuration file created at: %s", cmConfig.ConfigPath)

				return nil
			},
		},
		{
			Name:        "save",
			Aliases:     []string{"s"},
			Description: "Save the current configuration",
			Before:      beforeLoadConfig,
			Action: func(ctx *cli.Context) error {
				return nil
			},
		},
		{
			Name:        "load",
			Aliases:     []string{"l"},
			Description: "Load the onto your configuration",
			Before:      beforeLoadConfig,
			Action: func(ctx *cli.Context) error {
				return nil
			},
		},
	},
}
var (
	logger                      = cmLog.New(cmLog.Info)
	config *cmConfig.Definition = new(cmConfig.Definition)
)

func main() {
	if err := app.Run(os.Args); err != nil {
		logger.Error(err.Error())
	}
}

func beforeLoadConfig(ctx *cli.Context) error {
	// * Skip loading configuration when configuration
	if err := config.Load(logger); err != nil {
		return err
	}

	if config.LogLevel != "" {
		logger.Lvl = cmLog.LevelFromString(config.LogLevel)
	}

	return nil
}
