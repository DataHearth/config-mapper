package cmd

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	mapper "github.com/datahearth/config-mapper/internal"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	errLogger = log.New(os.Stderr, "", 0)
	logger    = log.New(os.Stderr, "", 0)
)

var rootCmd = &cobra.Command{
	Use:   "config-mapper",
	Short: "Manage your systems configuration",
	Long: `config-mapper aims to help you manage your configurations between systems
		with a single configuration file.`,
	Version: "v0.1.0.beta0",
}
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize your configuration folder",
	Long: `Initialize will retrieve your configuration folder from the source location and
		copy it into the destination field`,
	Run: func(cmd *cobra.Command, args []string) {
		var config mapper.Configuration

		if err := viper.Unmarshal(&config); err != nil {
			errLogger.Printf(pterm.Red(fmt.Sprintf("failed to decode configuration: %v\n", err)))
			os.Exit(1)
		}

		logger.Println("initializing config-mapper folder from configuration...")

		if _, err := mapper.NewRepository(config.Storage.Git, config.Storage.Path); err != nil {
			errLogger.Printf(pterm.Red(fmt.Sprintf("failed to initialize folder: %v\n", err)))
			os.Exit(1)
		}

		logger.Printf("repository initialized at \"%v\"\n", viper.GetString("storage.location"))
	},
}
var loadCmd = &cobra.Command{
	Use:   "load",
	Short: "Load your configurations onto your system",
	Long: `Load your files, folders and package managers deps configurations onto your new
		onto your new system based on your configuration file`,
	Run: func(cmd *cobra.Command, args []string) {
		var config mapper.Configuration

		if err := viper.Unmarshal(&config); err != nil {
			errLogger.Printf(pterm.Red(fmt.Sprintf("failed to decode configuration: %v\n", err)))
			os.Exit(1)
		}

		el := mapper.NewElement([]mapper.ItemLocation{}, config.Storage.Path)

		if !viper.GetBool("load-disable-files") {
			el.AddItems(config.Files)
		}
		if !viper.GetBool("load-disable-folders") {
			el.AddItems(config.Folders)
		}

		if err := el.Action("load"); err != nil {
			errLogger.Printf(pterm.Red(err))
			os.Exit(1)
		}

		if !viper.GetBool("load-disable-pkgs") {
			if err := mapper.LoadPkgs(config.PackageManagers); err != nil {
				errLogger.Printf(pterm.Red(err))
				os.Exit(1)
			}
		}
	},
}
var saveCmd = &cobra.Command{
	Use:   "save",
	Short: "save your configurations onto your saved location",
	Long: `Save your files, folders and package managers deps configurations onto your
		 saved location based on your configuration file`,
	Run: func(cmd *cobra.Command, args []string) {
		var config mapper.Configuration
		if err := viper.Unmarshal(&config); err != nil {
			errLogger.Printf(pterm.Red(fmt.Sprintf("failed to decode configuration: %v\n", err)))
			os.Exit(1)
		}

		repo, err := mapper.NewRepository(config.Storage.Git, config.Storage.Path)
		if err != nil {
			errLogger.Printf(pterm.Red(fmt.Sprintf("failed to open repository at %s: %v\n", config.Storage.Path, err)))
			os.Exit(1)
		}
		el := mapper.NewElement([]mapper.ItemLocation{}, config.Storage.Path)

		if !viper.GetBool("save-disable-files") {
			el.AddItems(config.Files)
		}
		if !viper.GetBool("save-disable-folders") {
			el.AddItems(config.Folders)
		}

		if err := el.Action("save"); err != nil {
			errLogger.Printf(pterm.Red(err))
			os.Exit(1)
		}

		if !viper.GetBool("save-disable-pkgs") {
			if err := mapper.SavePkgs(config); err != nil {
				errLogger.Printf(pterm.Red(err))
				os.Exit(1)
			}
		}

		if viper.GetBool("push") {
			pterm.DefaultSection.Println("Pushing items")

			s, _ := pterm.DefaultSpinner.WithShowTimer(true).WithRemoveWhenDone(false).Start("Pushing changes to remote repository")

			if err := repo.PushChanges(viper.GetString("message")); err != nil {
				errLogger.Printf(pterm.Red(fmt.Sprintf("failed to push changes to repository: %v\n", err)))
				os.Exit(1)
			}

			s.Stop()
			s.Success("Changes pushed")
		}
	},
}

func init() {
	cobra.OnInitialize(mapper.InitConfig)

	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(loadCmd)
	rootCmd.AddCommand(saveCmd)
	rootCmd.PersistentFlags().StringP("configuration-file", "c", "", "location of configuration file")
	viper.BindPFlag("configuration-file", rootCmd.PersistentFlags().Lookup("configuration-file"))

	loadCmd.PersistentFlags().Bool("disable-files", false, "files will be ignored")
	loadCmd.PersistentFlags().Bool("disable-folders", false, "folders will be ignored")
	loadCmd.PersistentFlags().Bool("disable-pkgs", false, "package managers will be ignored")
	viper.BindPFlag("load-disable-files", loadCmd.PersistentFlags().Lookup("disable-files"))
	viper.BindPFlag("load-disable-folders", loadCmd.PersistentFlags().Lookup("disable-folders"))
	viper.BindPFlag("load-disable-pkgs", loadCmd.PersistentFlags().Lookup("disable-pkgs"))

	saveCmd.PersistentFlags().Bool("disable-files", false, "files will be ignored")
	saveCmd.PersistentFlags().Bool("disable-folders", false, "folders will be ignored")
	saveCmd.PersistentFlags().Bool("disable-pkgs", false, "package managers will be ignored")
	saveCmd.Flags().BoolP("push", "p", false, "new configurations will be committed and pushed")
	saveCmd.Flags().StringP("message", "m", strconv.FormatInt(time.Now().Unix(), 10), "combined with --push to set a commit message")
	viper.BindPFlag("save-disable-files", saveCmd.PersistentFlags().Lookup("disable-files"))
	viper.BindPFlag("save-disable-folders", saveCmd.PersistentFlags().Lookup("disable-folders"))
	viper.BindPFlag("save-disable-pkgs", saveCmd.PersistentFlags().Lookup("disable-pkgs"))
	viper.BindPFlag("push", saveCmd.Flags().Lookup("push"))
	viper.BindPFlag("message", saveCmd.Flags().Lookup("message"))
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		errLogger.Printf("an error occured while running command: %v\n", err)
		os.Exit(1)
	}
}
