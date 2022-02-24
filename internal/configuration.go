package mapper

import (
	"fmt"
	"os"

	"github.com/pterm/pterm"
	"github.com/spf13/viper"
)

type Configuration struct {
	Storage         Storage     `yaml:"storage"`
	Files           []string    `yaml:"files"`
	Folders         []string    `yaml:"folders"`
	PackageManagers PkgManagers `yaml:"package-managers"`
}

type Storage struct {
	Location string `yaml:"location"`
	Git      Git    `yaml:"git"`
}

type Git struct {
	Username   string `yaml:"username"`
	Password   string `yaml:"password"`
	Repository string `yaml:"repository"`
}

type PkgManagers struct {
	InstallationOrder []string    `yaml:"installation-order"`
	Homebrew          []string    `yaml:"homebrew"`
	AptGet            interface{} `yaml:"apt-get"`
}

func InitConfig() {
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
	viper.SetConfigName(".config-mapper")

	viper.SetDefault("storage.location", fmt.Sprintf("%s/config-mapper", os.TempDir()))
	viper.SetDefault("package-managers.installation-order", []string{"apt", "homebrew"})

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			errLogger.Println(pterm.Red(err))
		} else {
			errLogger.Printf(pterm.Red(fmt.Sprintf("failed to read config: %v\n", err)))
		}

		os.Exit(1)
	}
}
