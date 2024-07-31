package mapper

import (
	"fmt"
	"os"

	"github.com/pterm/pterm"
	"github.com/spf13/viper"
)

type Configuration struct {
	Storage         Storage        `mapstructure:"storage"`
	Files           []ItemLocation `mapstructure:"files"`
	Folders         []ItemLocation `mapstructure:"folders"`
	PackageManagers PkgManagers    `mapstructure:"package-managers"`
}

type ItemLocation struct {
	Darwin string `mapstructure:"darwin"`
	Linux  string `mapstructure:"linux"`
}

type Storage struct {
	Location string `mapstructure:"location"`
	Git      Git    `mapstructure:"git"`
}

type Git struct {
	SSH        Ssh       `mapstructure:"ssh"`
	BasicAuth  BasicAuth `mapstructure:"basic-auth"`
	Repository string    `mapstructure:"repository"`
}

type BasicAuth struct {
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

type Ssh struct {
	Passphrase string `mapstructure:"passphrase"`
	PrivateKey string `mapstructure:"private-key"`
}

type PkgManagers struct {
	InstallationOrder []string `mapstructure:"installation-order"`
	Homebrew          []string `mapstructure:"homebrew"`
	Aptitude          []string `mapstructure:"apt-get"`
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
