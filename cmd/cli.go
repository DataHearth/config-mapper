package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/datahearth/config-mapper/internal/git"
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
}
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize your configuration folder",
	Long: `Initialize will retrieve your configuration folder from the source location and
		copy it into the destination field`,
	Run: func(cmd *cobra.Command, args []string) {
		logger.Println("initializing config-mapper folder from configuration...")
		if _, err := git.OpenGitRepo(); err != nil {
			errLogger.Printf("failed to initialize folder: %v\n", err)
			os.Exit(1)
		}
		logger.Printf("repository initialized at \"%v\"\n", viper.GetString("storage.location"))
	},
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.AddCommand(initCmd)
}

func initConfig() {
	h, err := os.UserHomeDir()
	if err != nil {
		errLogger.Printf("can't get home directory through $HOME variable: %v\n", err)
		os.Exit(1)
	}

	if c := os.Getenv("CONFIG_MAPPER_CFG"); c != "" {
		viper.AddConfigPath(c)
	} else {
		viper.AddConfigPath(h)
	}

	viper.SetConfigType("yml")
	viper.SetConfigName("config-mapper")

	viper.SetDefault("storage.location", fmt.Sprintf("%s/config-mapper", os.TempDir()))

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			errLogger.Println("configuration file not found", err)
		} else {
			errLogger.Printf("failed to read config: %v\n", err)
		}

		os.Exit(1)
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		errLogger.Printf("an error occured while running command: %v\n", err)
		os.Exit(1)
	}
}
