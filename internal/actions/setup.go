package actions

import (
	"os"
	"path/filepath"

	cmConfig "gitea.antoine-langlois.net/datahearth/config-mapper/internal/config"
	"github.com/google/uuid"
)

// Setup is the action to setup the configuration file
func Setup() error {
	cfg := cmConfig.Definition{
		UUID:     cmConfig.UUIDv4{UUID: uuid.New()},
		Path:     "$HOME/.local/state/config-mapper",
		LogLevel: "info",
		Storage:  cmConfig.Storage{},
		Items:    []string{},
	}

	if err := os.MkdirAll(filepath.Dir(os.ExpandEnv(cmConfig.ConfigPath)), 0755); err != nil {
		return err
	}

	return cfg.Write()
}
