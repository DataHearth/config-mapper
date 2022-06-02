package configuration

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/viper"
)

var errLogger = log.New(os.Stderr, "", 0)

func InitConfig() {
	h, err := os.UserHomeDir()
	if err != nil {
		errLogger.Fatalln(err)
	}

	if c := viper.GetString("configuration-file"); c != "" {
		if strings.Contains(c, "ssh://") {
			viper.AddConfigPath(h)

			viper.SetConfigType("yml")
			viper.SetConfigName(".config-mapper")

			if err := loadConfigSSH(c); err != nil {
				errLogger.Fatalln(err)
			}
			return
		}

		if strings.Contains(c, ".yml") {
			viper.AddConfigPath(path.Dir(c))
		} else {
			viper.AddConfigPath(c)
		}
	}
	if c := os.Getenv("CONFIG_MAPPER_CFG"); c != "" {
		if strings.Contains(c, ".yml") {
			viper.AddConfigPath(path.Dir(c))
		} else {
			viper.AddConfigPath(c)
		}
	}
	viper.AddConfigPath(h)

	viper.SetConfigType("yml")
	viper.SetConfigName(".config-mapper")

	viper.SetDefault("storage.location", fmt.Sprintf("%s/config-mapper", os.TempDir()))
	viper.SetDefault("package-managers.installation-order", []string{"apt", "homebrew"})

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			color.Error.Write([]byte(color.RedString("no configuration file found: %v\n", err)))
		} else {
			color.Error.Write([]byte(color.RedString("failed to read config: %v\n", err)))
		}

		os.Exit(1)
	}

	viper.Set("configuration-file", viper.ConfigFileUsed())
}
