package configuration

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/spf13/viper"
)

func InitConfig() {
	h, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	if c := viper.GetString("configuration-file"); c != "" {
		if strings.Contains(c, "ssh://") {
			viper.AddConfigPath(h)

			viper.SetConfigType("yml")
			viper.SetConfigName(".config-mapper")

			if err := loadConfigSSH(c); err != nil {
				log.Fatal(err)
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

	if err := viper.ReadInConfig(); err != nil {
		var errMsg string
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			errMsg = "no configuration file found"
		} else {
			errMsg = "failed to read config"
		}

		log.Fatal(errMsg, "err", err)
	}

	viper.Set("configuration-file", viper.ConfigFileUsed())
}
