package config

import (
	"io"
	"os"

	cmLog "gitea.antoine-langlois.net/datahearth/config-mapper/internal/logger"
	"github.com/google/uuid"
	"gopkg.in/yaml.v3"
)

type Definition struct {
	UUID     uuid.UUID `yaml:"uuid"`
	Path     string    `yaml:"path,omitempty"`
	LogLevel string    `yaml:"log-level,omitempty"`
	Storage  Storage   `yaml:"storage"`
	Items    []string  `yaml:"items"`
}

func Load(logger *cmLog.Logger) (*Definition, error) {
	env := "$HOME/.config/config-mapper.yml"
	if v := os.Getenv("CFG_MAPPER_CONFIG_PATH"); v != "" {
		env = v
	}

	logger.Debug("Loading config from %s", env)
	f, err := os.Open(env)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	cfg := Definition{}
	d, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(d, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
