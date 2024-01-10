package config

import (
	"io"
	"os"

	cmLog "gitea.antoine-langlois.net/datahearth/config-mapper/internal/logger"
	"github.com/google/uuid"
	"gopkg.in/yaml.v3"
)

const ConfigPath = "$HOME/.config/config-mapper.yml"

type UUIDv4 struct {
	uuid.UUID
}

func (u *UUIDv4) UnmarshalYAML(v *yaml.Node) error {
	var uuidStr string
	if err := v.Decode(&uuidStr); err != nil {
		return err
	}

	uuid, err := uuid.Parse(uuidStr)
	if err != nil {
		return err
	}

	u.UUID = uuid

	return nil
}

type Definition struct {
	// unique system identifier
	UUID UUIDv4 `yaml:"uuid"`
	// directory location where the configuration is stored
	Path string `yaml:"path,omitempty"`
	// override log level
	LogLevel string `yaml:"log-level,omitempty"`
	// storage providers configuration
	Storage Storage `yaml:"storage"`
	// registered items
	Items []string `yaml:"items"`
}

func (d *Definition) Load(logger *cmLog.Logger) error {
	env := "$HOME/.config/config-mapper.yml"
	if v := os.Getenv("CFG_MAPPER_CONFIG_PATH"); v != "" {
		env = v
	}

	logger.Debug("Loading config from %s", env)
	f, err := os.Open(env)
	if err != nil {
		return err
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		return err
	}

	if err := yaml.Unmarshal(data, d); err != nil {
		return err
	}

	return nil
}

func (d *Definition) Write() error {
	b, err := yaml.Marshal(d)
	if err != nil {
		return err
	}

	return os.WriteFile(os.ExpandEnv(ConfigPath), b, 0644)
}
