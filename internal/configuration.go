package mapper

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/pterm/pterm"
	"github.com/spf13/viper"
)

type Configuration struct {
	Storage         Storage        `mapstructure:"storage" yaml:"storage"`
	Files           []ItemLocation `mapstructure:"files" yaml:"files"`
	Folders         []ItemLocation `mapstructure:"folders" yaml:"folders"`
	PackageManagers PkgManagers    `mapstructure:"package-managers" yaml:"package-managers"`
}

type ItemLocation struct {
	Darwin string `mapstructure:"darwin" yaml:"darwin"`
	Linux  string `mapstructure:"linux" yaml:"linux"`
}

type Storage struct {
	Location string `mapstructure:"location" yaml:"location"`
	Git      Git    `mapstructure:"git" yaml:"git"`
}

type Git struct {
	Repository string    `mapstructure:"repository" yaml:"repository"`
	BasicAuth  BasicAuth `mapstructure:"basic-auth" yaml:"basic-auth"`
	SSH        Ssh       `mapstructure:"ssh" yaml:"ssh"`
}

type BasicAuth struct {
	Username string `mapstructure:"username" yaml:"username"`
	Password string `mapstructure:"password" yaml:"password"`
}

type Ssh struct {
	PrivateKey string `mapstructure:"private-key" yaml:"private-key"`
	Passphrase string `mapstructure:"passphrase" yaml:"passphrase"`
}

type PkgManagers struct {
	InstallationOrder []string `mapstructure:"installation-order" yaml:"installation-order"`
	Homebrew          []string `mapstructure:"homebrew" yaml:"homebrew"`
	Aptitude          []string `mapstructure:"apt-get" yaml:"apt-get"`
}

func InitConfig() {
	h, err := os.UserHomeDir()
	if err != nil {
		errLogger.Printf("can't get home directory through $HOME variable: %v\n", err)
		os.Exit(1)
	}

	if c := viper.GetString("configuration-file"); c != "" {
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
			errLogger.Println(pterm.Red(err))
		} else {
			errLogger.Printf(pterm.Red(fmt.Sprintf("failed to read config: %v\n", err)))
		}

		os.Exit(1)
	}

	viper.Set("configuration-file", viper.ConfigFileUsed())
}
