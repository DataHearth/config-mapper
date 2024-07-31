package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var errLogger = log.New(os.Stderr, "", 0)

var rootCmd = &cobra.Command{
	Use:   "config-mapper",
	Short: "Manage your systems configuration",
	Long: `config-mapper aims to help you manage your configurations between systems
						with a single configuration file.`,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	h, err := os.UserHomeDir()
	if err != nil {
		errLogger.Printf("can't get home directory through $HOME variable: %v\n", err)
		os.Exit(1)
	}

	if c := os.Getenv("CONFIG_MAPPER"); c != "" {
		viper.AddConfigPath(c)
	} else {
		viper.AddConfigPath(h)
	}

	viper.SetConfigType("yml")
	viper.SetConfigName("config-mapper")
	if err := viper.ReadInConfig(); err != nil {
		errLogger.Printf("failed to read config: %v\n", err)
		os.Exit(1)
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		errLogger.Printf("an error occured while running command: %v\n", err)
		os.Exit(1)
	}
}
