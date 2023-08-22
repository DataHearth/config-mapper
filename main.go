package main

import (
	"strconv"
	"time"

	"gitea.antoine-langlois.net/datahearth/config-mapper/internal/logging"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "config-mapper",
	Short: "Manage your systems configuration",
	Long: `config-mapper aims to help you manage your configurations between systems
		with a single configuration file.`,
	Version: "v0.6.2",
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
		logrus.Fatal("install command not implemented yet")
	},
}

func init() {
	logrus.SetFormatter(new(logging.LoggerFormatter))

	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(loadCmd)
	rootCmd.AddCommand(saveCmd)
	rootCmd.AddCommand(installCmd)

	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "STDOUT will be more verbose")
	rootCmd.PersistentFlags().StringP("configuration-file", "c", "", "location of configuration file")
	rootCmd.PersistentFlags().String("ssh-user", "", "SSH username to retrieve configuration file")
	rootCmd.PersistentFlags().String("ssh-password", "", "SSH password to retrieve configuration file")
	rootCmd.PersistentFlags().String("ssh-key", "", "SSH key to retrieve configuration file (if a passphrase is needed, use the \"CONFIG_MAPPER_PASS\" env variable")

	loadCmd.Flags().Bool("disable-files", false, "files will be ignored")
	loadCmd.Flags().Bool("disable-folders", false, "folders will be ignored")
	loadCmd.Flags().Bool("pkgs", false, "packages will be installed")
	loadCmd.Flags().StringSlice("exclude-pkg-managers", []string{}, "package managers to exclude (comma separated)")

	saveCmd.Flags().Bool("disable-files", false, "files will be ignored")
	saveCmd.Flags().Bool("disable-folders", false, "folders will be ignored")
	saveCmd.Flags().BoolP("push", "p", false, "new configurations will be committed and pushed")
	saveCmd.Flags().StringP("message", "m", strconv.FormatInt(time.Now().Unix(), 10), "combined with --push to set a commit message")
	saveCmd.Flags().Bool("disable-index", false, "configuration index will not be updated")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		logrus.Fatal("an error occured while running command", "err", err)
	}
}

func save(cmd *cobra.Command, args []string) {
}

func load(cmd *cobra.Command, args []string) {
}

// func save(cmd *cobra.Command, args []string) {
// 	var c configuration.Configuration
// 	if err := viper.Unmarshal(&c); err != nil {
// 		logrus.Fatal("failed to decode configuration", "err", err)
// 	}

// 	indexer, err := mapper.NewIndexer(c.Storage.Path)
// 	if err != nil {
// 		logrus.Fatal("failed to open the indexer", "err", err)
// 	}

// 	r, err := git.NewRepository(c.Storage.Git, c.Storage.Path)
// 	if err != nil {
// 		logrus.Fatal("failed to open repository", "path", c.Storage.Path, "err", err)
// 	}

// 	el := mapper.NewItemsActions(nil, c.Storage.Path, r, indexer)

// 	if !viper.GetBool("save-disable-files") {
// 		el.AddItems(c.Files)
// 	}
// 	if !viper.GetBool("save-disable-folders") {
// 		el.AddItems(c.Folders)
// 	}

// 	el.Action("save")

// 	if err := el.CleanUp(indexer.RemovedLines()); err != nil {
// 		logrus.Fatal("failed to clean repository", "err", err)
// 	}

// 	if viper.GetBool("push") {
// 		logrus.Info("pushing changes...")

// 		if err := r.PushChanges(viper.GetString("message"), indexer.Lines(), indexer.RemovedLines()); err != nil {
// 			logrus.Fatal("failed to push changes to repository", "err", err)
// 		}
// 	}
// }

// func load(cmd *cobra.Command, args []string) {
// 	var c configuration.Configuration
// 	if err := viper.Unmarshal(&c); err != nil {
// 		logrus.Fatal("failed to decode configuration", "err", err)
// 	}

// 	i, err := mapper.NewIndexer(c.Storage.Path)
// 	if err != nil {
// 		logrus.Fatal("failed to open the indexer", "err", err)
// 	}

// 	r, err := git.NewRepository(c.Storage.Git, c.Storage.Path)
// 	if err != nil {
// 		logrus.Fatal("failed to open repository", "path", c.Storage.Path, "err", err)
// 	}

// 	el := mapper.NewItemsActions(nil, c.Storage.Path, r, i)

// 	if !viper.GetBool("load-disable-files") {
// 		el.AddItems(c.Files)
// 	}
// 	if !viper.GetBool("load-disable-folders") {
// 		el.AddItems(c.Folders)
// 	}

// 	el.Action("load")

// 	if viper.GetBool("load-enable-pkgs") {
// 		if err := mapper.InstallPackages(c.PackageManagers); err != nil {
// 			logrus.Fatal(err)
// 		}
// 	}
// }

func initCommand(cmd *cobra.Command, args []string) {
}
