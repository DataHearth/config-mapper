package cmd

import (
	"fmt"
	"log"
	"os"

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

		if _, err := mapper.OpenGitRepo(config.Storage.Git, config.Storage.Location); err != nil {
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

		if !viper.GetBool("disable-files") {
			if err := mapper.LoadFiles(config.Files, config.Storage.Location); err != nil {
				errLogger.Printf(pterm.Red(fmt.Sprintf("error while loading files: %v\n", err)))
				os.Exit(1)
			}
		}
		if !viper.GetBool("disable-pkgs") {
			if err := mapper.LoadPkgs(config.PackageManagers); err != nil {
				errLogger.Printf(pterm.Red(fmt.Sprintf("error while installing packages: %v\n", err)))
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

		if !viper.GetBool("disable-files") {
			if err := mapper.SaveFiles(config.Files, config.Storage.Location); err != nil {
				errLogger.Printf(pterm.Red(fmt.Sprintf("error while loading files: %v\n", err)))
				os.Exit(1)
			}
		}
		if !viper.GetBool("disable-folders") {
			if err := mapper.SaveFolders(config.Folders, config.Storage.Location); err != nil {
				errLogger.Printf(pterm.Red(fmt.Sprintf("error while loading files: %v\n", err)))
				os.Exit(1)
			}
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
	viper.BindPFlag("disable-files", loadCmd.PersistentFlags().Lookup("disable-files"))
	viper.BindPFlag("disable-folders", loadCmd.PersistentFlags().Lookup("disable-folders"))
	viper.BindPFlag("disable-pkgs", loadCmd.PersistentFlags().Lookup("disable-pkgs"))

	saveCmd.PersistentFlags().Bool("disable-files", false, "files will be ignored")
	saveCmd.PersistentFlags().Bool("disable-folders", false, "folders will be ignored")
	saveCmd.PersistentFlags().Bool("disable-pkgs", false, "package managers will be ignored")
	viper.BindPFlag("disable-files", saveCmd.PersistentFlags().Lookup("disable-files"))
	viper.BindPFlag("disable-folders", saveCmd.PersistentFlags().Lookup("disable-folders"))
	viper.BindPFlag("disable-pkgs", saveCmd.PersistentFlags().Lookup("disable-pkgs"))
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		errLogger.Printf("an error occured while running command: %v\n", err)
		os.Exit(1)
	}
}
