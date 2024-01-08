package main

import (
	"os"

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
	Before: func(ctx *cli.Context) error {
		var err error
		config, err = cmConfig.Load(logger)
		if err != nil {
			return err
		}

		if config.LogLevel != "" {
			logger.Lvl = cmLog.LevelFromString(config.LogLevel)
		}

		return nil
	},
	Commands: []*cli.Command{},
}
var (
	logger = cmLog.New(cmLog.Info)
	config *cmConfig.Definition
)

func main() {
	if err := app.Run(os.Args); err != nil {
		logger.Error(err.Error())
	}
}
