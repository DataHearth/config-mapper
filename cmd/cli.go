package cmd

import (
	"log"
	"os"
	"strconv"
	"time"

	mapper "github.com/datahearth/config-mapper/internal"
	"github.com/datahearth/config-mapper/internal/configuration"
	"github.com/datahearth/config-mapper/internal/git"
	"github.com/fatih/color"
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
	Version: "v0.4.0",
}
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize your configuration folder",
	Long: `Initialize will retrieve your configuration folder from the source location and
		copy it into the destination field`,
	Run: initCommand,
}
var loadCmd = &cobra.Command{
	Use:   "load",
	Short: "Load your configurations onto your system",
	Long: `Load your files, folders and package managers deps configurations onto your new
		onto your new system based on your configuration file`,
	Run: load,
}
var saveCmd = &cobra.Command{
	Use:   "save",
	Short: "save your configurations onto your saved location",
	Long: `Save your files, folders and package managers deps configurations onto your
		 saved location based on your configuration file`,
	Run: save,
}

func init() {
	cobra.OnInitialize(configuration.InitConfig)

	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(loadCmd)
	rootCmd.AddCommand(saveCmd)
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "STDOUT will be more verbose")
	rootCmd.PersistentFlags().StringP("configuration-file", "c", "", "location of configuration file")
	rootCmd.PersistentFlags().String("ssh-user", "", "SSH username to retrieve configuration file")
	rootCmd.PersistentFlags().String("ssh-password", "", "SSH password to retrieve configuration file")
	rootCmd.PersistentFlags().String("ssh-key", "", "SSH key to retrieve configuration file (if a passphrase is needed, use the \"CONFIG_MAPPER_PASS\" env variable")
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
	viper.BindPFlag("configuration-file", rootCmd.PersistentFlags().Lookup("configuration-file"))
	viper.BindPFlag("ssh-user", rootCmd.PersistentFlags().Lookup("ssh-user"))
	viper.BindPFlag("ssh-password", rootCmd.PersistentFlags().Lookup("ssh-password"))
	viper.BindPFlag("ssh-key", rootCmd.PersistentFlags().Lookup("ssh-key"))

	loadCmd.Flags().Bool("disable-files", false, "files will be ignored")
	loadCmd.Flags().Bool("disable-folders", false, "folders will be ignored")
	loadCmd.Flags().Bool("pkgs", false, "packages will be installed")
	viper.BindPFlag("load-disable-files", loadCmd.Flags().Lookup("disable-files"))
	viper.BindPFlag("load-disable-folders", loadCmd.Flags().Lookup("disable-folders"))
	viper.BindPFlag("load-enable-pkgs", loadCmd.Flags().Lookup("pkgs"))

	saveCmd.Flags().Bool("disable-files", false, "files will be ignored")
	saveCmd.Flags().Bool("disable-folders", false, "folders will be ignored")
	saveCmd.Flags().Bool("pkgs", false, "packages will be saved")
	saveCmd.Flags().BoolP("push", "p", false, "new configurations will be committed and pushed")
	saveCmd.Flags().StringP("message", "m", strconv.FormatInt(time.Now().Unix(), 10), "combined with --push to set a commit message")
	saveCmd.Flags().Bool("disable-index", false, "configuration index will not be updated")
	viper.BindPFlag("save-disable-files", saveCmd.Flags().Lookup("disable-files"))
	viper.BindPFlag("save-disable-folders", saveCmd.Flags().Lookup("disable-folders"))
	viper.BindPFlag("save-enable-pkgs", saveCmd.Flags().Lookup("pkgs"))
	viper.BindPFlag("push", saveCmd.Flags().Lookup("push"))
	viper.BindPFlag("disable-index-update", saveCmd.Flags().Lookup("disable-index"))
	viper.BindPFlag("message", saveCmd.Flags().Lookup("message"))
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		errLogger.Printf("an error occured while running command: %v\n", err)
		os.Exit(1)
	}
}

func save(cmd *cobra.Command, args []string) {
	var c configuration.Configuration
	if err := viper.Unmarshal(&c); err != nil {
		mapper.PrintError("failed to decode configuration: %v\n", err)
		os.Exit(1)
	}

	indexer, err := mapper.NewIndexer(c.Storage.Path)
	if err != nil {
		mapper.PrintError("failed to open the indexer: %v\n", err)
		os.Exit(1)
	}

	r, err := git.NewRepository(c.Storage.Git, c.Storage.Path)
	if err != nil {
		mapper.PrintError("failed to open repository at %s: %v\n", c.Storage.Path, err)
		os.Exit(1)
	}

	el := mapper.NewItemsActions(nil, c.Storage.Path, r, indexer)

	if !viper.GetBool("save-disable-files") {
		el.AddItems(c.Files)
	}
	if !viper.GetBool("save-disable-folders") {
		el.AddItems(c.Folders)
	}

	el.Action("save")

	if viper.GetBool("save-enable-pkgs") {
		if err := mapper.SavePkgs(c); err != nil {
			mapper.PrintError(err.Error())
			os.Exit(1)
		}
	}

	if err := el.CleanUp(indexer.RemovedLines()); err != nil {
		mapper.PrintError("failed to clean repository: %v\n", err)
		os.Exit(1)
	}

	if viper.GetBool("push") {
		color.Blue("# Pushing items")

		if err := r.PushChanges(viper.GetString("message"), indexer.Lines(), indexer.RemovedLines()); err != nil {
			mapper.PrintError("failed to push changes to repository: %v\n", err)
			os.Exit(1)
		}

		color.Green("Items pushed")
	}
}

func load(cmd *cobra.Command, args []string) {
	var c configuration.Configuration
	if err := viper.Unmarshal(&c); err != nil {
		mapper.PrintError("failed to decode configuration: %v\n", err)
		os.Exit(1)
	}

	i, err := mapper.NewIndexer(c.Storage.Path)
	if err != nil {
		mapper.PrintError("failed to open the indexer: %v\n", err)
		os.Exit(1)
	}

	r, err := git.NewRepository(c.Storage.Git, c.Storage.Path)
	if err != nil {
		mapper.PrintError("failed to open repository at %s: %v\n", c.Storage.Path, err)
		os.Exit(1)
	}

	el := mapper.NewItemsActions(nil, c.Storage.Path, r, i)

	if !viper.GetBool("load-disable-files") {
		el.AddItems(c.Files)
	}
	if !viper.GetBool("load-disable-folders") {
		el.AddItems(c.Folders)
	}

	el.Action("load")

	if viper.GetBool("load-enable-pkgs") {
		if err := mapper.LoadPkgs(c.PackageManagers); err != nil {
			mapper.PrintError(err.Error())
			os.Exit(1)
		}
	}
}

func initCommand(cmd *cobra.Command, args []string) {
	var c configuration.Configuration
	if err := viper.Unmarshal(&c); err != nil {
		mapper.PrintError("failed to decode configuration: %v\n", err)
		os.Exit(1)
	}

	logger.Println("initializing config-mapper folder from configuration...")

	if _, err := git.NewRepository(c.Storage.Git, c.Storage.Path); err != nil {
		mapper.PrintError("failed to initialize folder: %v\n", err)
		os.Exit(1)
	}

	logger.Printf("repository initialized at \"%v\"\n", viper.GetString("storage.location"))
}
