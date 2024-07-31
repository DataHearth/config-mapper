package cmd

import (
	"strconv"
	"time"

	mapper "gitea.antoine-langlois.net/datahearth/config-mapper/internal"
	"gitea.antoine-langlois.net/datahearth/config-mapper/internal/configuration"
	"gitea.antoine-langlois.net/datahearth/config-mapper/internal/git"
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "config-mapper",
	Short: "Manage your systems configuration",
	Long: `config-mapper aims to help you manage your configurations between systems
		with a single configuration file.`,
	Version: "v0.6.1",
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
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "install additional tools",
	Long:  `install additional tools like package managers, programming languages, etc.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Fatal("install command not implemented yet")
	},
}

func init() {
	cobra.OnInitialize(configuration.InitConfig)

	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(loadCmd)
	rootCmd.AddCommand(saveCmd)
	rootCmd.AddCommand(installCmd)

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
	loadCmd.Flags().StringSlice("exclude-pkg-managers", []string{}, "package managers to exclude (comma separated)")
	viper.BindPFlag("load-disable-files", loadCmd.Flags().Lookup("disable-files"))
	viper.BindPFlag("load-disable-folders", loadCmd.Flags().Lookup("disable-folders"))
	viper.BindPFlag("load-enable-pkgs", loadCmd.Flags().Lookup("pkgs"))
	viper.BindPFlag("exclude-pkg-managers", loadCmd.Flags().Lookup("exclude-pkg-managers"))

	saveCmd.Flags().Bool("disable-files", false, "files will be ignored")
	saveCmd.Flags().Bool("disable-folders", false, "folders will be ignored")
	saveCmd.Flags().BoolP("push", "p", false, "new configurations will be committed and pushed")
	saveCmd.Flags().StringP("message", "m", strconv.FormatInt(time.Now().Unix(), 10), "combined with --push to set a commit message")
	saveCmd.Flags().Bool("disable-index", false, "configuration index will not be updated")
	viper.BindPFlag("save-disable-files", saveCmd.Flags().Lookup("disable-files"))
	viper.BindPFlag("save-disable-folders", saveCmd.Flags().Lookup("disable-folders"))
	viper.BindPFlag("push", saveCmd.Flags().Lookup("push"))
	viper.BindPFlag("disable-index-update", saveCmd.Flags().Lookup("disable-index"))
	viper.BindPFlag("message", saveCmd.Flags().Lookup("message"))
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal("an error occured while running command", "err", err)
	}
}

func save(cmd *cobra.Command, args []string) {
	var c configuration.Configuration
	if err := viper.Unmarshal(&c); err != nil {
		log.Fatal("failed to decode configuration", "err", err)
	}

	indexer, err := mapper.NewIndexer(c.Storage.Path)
	if err != nil {
		log.Fatal("failed to open the indexer", "err", err)
	}

	r, err := git.NewRepository(c.Storage.Git, c.Storage.Path)
	if err != nil {
		log.Fatal("failed to open repository", "path", c.Storage.Path, "err", err)
	}

	el := mapper.NewItemsActions(nil, c.Storage.Path, r, indexer)

	if !viper.GetBool("save-disable-files") {
		el.AddItems(c.Files)
	}
	if !viper.GetBool("save-disable-folders") {
		el.AddItems(c.Folders)
	}

	el.Action("save")

	if err := el.CleanUp(indexer.RemovedLines()); err != nil {
		log.Fatal("failed to clean repository", "err", err)
	}

	if viper.GetBool("push") {
		log.Info("pushing changes...")

		if err := r.PushChanges(viper.GetString("message"), indexer.Lines(), indexer.RemovedLines()); err != nil {
			log.Fatal("failed to push changes to repository", "err", err)
		}
	}
}

func load(cmd *cobra.Command, args []string) {
	var c configuration.Configuration
	if err := viper.Unmarshal(&c); err != nil {
		log.Fatal("failed to decode configuration", "err", err)
	}

	i, err := mapper.NewIndexer(c.Storage.Path)
	if err != nil {
		log.Fatal("failed to open the indexer", "err", err)
	}

	r, err := git.NewRepository(c.Storage.Git, c.Storage.Path)
	if err != nil {
		log.Fatal("failed to open repository", "path", c.Storage.Path, "err", err)
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
		if err := mapper.InstallPackages(c.PackageManagers); err != nil {
			log.Fatal(err)
		}
	}
}

func initCommand(cmd *cobra.Command, args []string) {
	var c configuration.Configuration
	if err := viper.Unmarshal(&c); err != nil {
		log.Fatal("failed to decode configuration", "err", err)
	}

	log.Info("initializing config-mapper folder from configuration...")

	if _, err := git.NewRepository(c.Storage.Git, c.Storage.Path); err != nil {
		log.Fatal("failed to initialize folder", "err", err)
	}

	log.Info("repository initialized", "path", viper.GetString("storage.location"))
}
