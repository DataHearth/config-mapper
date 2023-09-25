package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"gitea.antoine-langlois.net/datahearth/config-mapper/internal"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"golang.org/x/exp/slices"
	"gopkg.in/yaml.v3"
)

var (
	configuration internal.Configuration
	logFormatter  = new(internal.LoggerFormatter)
	app           = &cli.App{
		Version:     "v0.6.2",
		Name:        "config-mapper",
		Description: "Manage your systems configuration",
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
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Usage:   "location of configuration file",
				Value:   "$HOME/.config-mapper.yml",
				Action: func(ctx *cli.Context, s string) error {
					path, err := internal.ResolvePath(s)
					if err != nil {
						return err
					}
					stat, err := os.Stat(path)
					if err != nil {
						return err
					}

					if stat.IsDir() {
						return fmt.Errorf("--config must be a file. Found directory")
					}

					if !slices.Contains([]string{".yml", ".yaml"}, filepath.Ext(path)) {
						return fmt.Errorf("--config must be a \".yml|yaml\" file")
					}

					return nil
				},
			},
		},
		Commands: []*cli.Command{
			{
				Name:    "init",
				Aliases: []string{"i"},
				Usage:   "Initialize your configuration folder",
				UsageText: `Initialize will retrieve your configuration folder from the source location and
					copy it into the destination field`,
				Action: initCommand,
				Before: beforeAction,
			},
			{
				Name:    "load",
				Aliases: []string{"l"},
				Usage:   "Load your configurations onto your system",
				UsageText: `Load your items and package managers dependencies onto your new
					system based on your configuration file`,
				Action: loadCommand,
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:     "items",
						Usage:    "Items will not be loaded",
						Aliases:  []string{"i"},
						Category: "Exclude",
					},
					&cli.BoolFlag{
						Name:    "packages",
						Usage:   "Packages will be installed",
						Aliases: []string{"p"},
					},
					&cli.StringSliceFlag{
						Name:     "package-managers",
						Usage:    "Exclude packages from specified package managers to be installed",
						Category: "Exclude",
						Aliases:  []string{"P"},
						Action: func(ctx *cli.Context, s []string) error {
							if !ctx.Bool("packages") {
								return fmt.Errorf("--packages is required to exclude package managers")
							}

							return nil
						},
					},
				},
				Before: beforeAction,
			},
			{
				Name:      "save",
				Aliases:   []string{"s"},
				Usage:     "save your configurations onto your saved location",
				UsageText: `Save your items onto your saved location based on your configuration file`,
				Action:    saveCommand,
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:    "push",
						Usage:   "Push changes to remote repository",
						Aliases: []string{"p"},
					},
					&cli.StringFlag{
						Name:    "message",
						Usage:   "Commit message when pushing repository",
						Aliases: []string{"m"},
						Action: func(ctx *cli.Context, s string) error {
							if !ctx.Bool("push") {
								return fmt.Errorf("--push option is required to set a message")
							}
							if len(strings.Trim(s, " ")) == 0 {
								return fmt.Errorf("message must contain at least one character")
							}

							return nil
						},
					},
				},
				Before: beforeAction,
			},
			{
				Name:  "install",
				Usage: "Install additional tools",
				UsageText: `Install additional tools like package managers (brew, nala), 
					programming language (golang, rust, python)`,
				ArgsUsage: "TOOLS",
				Action:    installCommand,
			},
			{
				Name:    "check",
				Aliases: []string{"c"},
				Usage:   "Check if remote has updates",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:  "pull",
						Usage: "Pull updates if any",
					},
				},
				Action: checkCommand,
				Before: beforeAction,
			},
		},
	}
)

func init() {
	logrus.SetFormatter(logFormatter)
}

func main() {
	if err := app.Run(os.Args); err != nil {
		logrus.Fatalln(err)
	}
}

func initCommand(ctx *cli.Context) error {
	if _, err := os.Stat(configuration.Path); err == nil {
		logrus.Warnf("configuration folder already exists at %s", configuration.Path)
		return nil
	}

	logrus.WithField("path", configuration.Path).Infoln("initializing configuration folder...")

	if _, err := internal.NewRepository(configuration.Storage.Git, configuration.Path, true); err != nil {
		return err
	}

	logrus.Infoln("Everything is ready to go!")

	return nil
}

func saveCommand(ctx *cli.Context) error {
	repo, err := internal.NewRepository(configuration.Storage.Git, configuration.Path, false)
	if err != nil {
		return err
	}

	if configuration.Storage.Git.Pull {
		if err := repo.Pull(); err != nil {
			return err
		}

		logrus.Debugln("Changes pulled")
	}

	return nil
}

func loadCommand(ctx *cli.Context) error {
	return nil
}

func installCommand(ctx *cli.Context) error {
	return nil
}

func checkCommand(ctx *cli.Context) error {
	repo, err := internal.NewRepository(configuration.Storage.Git, configuration.Path, false)
	if err != nil {
		return err
	}

	changes, err := repo.FetchChanges()
	if err != nil {
		return err
	}
	if changes {
		logrus.Infoln("Configuration has changed upstream")
		if ctx.Bool("pull") {
			logrus.Infoln("Pulling changes...")
			if err := repo.Pull(); err != nil {
				return err
			}
		}
	} else {
		logrus.Infoln("Configuration is up-to-date")
	}

	return nil
}

func beforeAction(ctx *cli.Context) error {
	configPath, err := internal.ResolvePath(ctx.String("config"))
	if err != nil {
		return err
	}

	f, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}

	if err := yaml.Unmarshal(f, &configuration); err != nil {
		return err
	}

	if configuration.Logging.Level != "" {
		lvl, err := logrus.ParseLevel(configuration.Logging.Level)
		if err != nil {
			return err
		}

		logrus.SetLevel(lvl)
	}
	if configuration.Logging.File != "" {
		loggingFile := configuration.Logging.File
		s, err := os.Stat(filepath.Dir(loggingFile))
		if err != nil {
			if !os.IsNotExist(err) {
				return err
			}

			if err := os.MkdirAll(filepath.Dir(loggingFile), 0755); err != nil {
				return err
			}
		}

		if !s.IsDir() {
			return fmt.Errorf("parent path segment of \"log-file\" field is a file")
		}

		f, err := os.OpenFile(configuration.Logging.File, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0755)
		if err != nil {
			return err
		}

		var out io.Writer
		if configuration.Logging.DisableConsole {
			out = f
		} else {
			out = io.MultiWriter(os.Stdout, f)
		}

		logrus.SetOutput(out)
	}

	if configuration.Logging.EnableTime {
		format := "02/01/2006 15:04:05"
		if configuration.Logging.TimeFormat != "" {
			format = configuration.Logging.TimeFormat
		}

		logFormatter.TimeFormat = format
	}

	return nil
}
