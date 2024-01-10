package actions

import (
	cmConfig "gitea.antoine-langlois.net/datahearth/config-mapper/internal/config"
	cmProviders "gitea.antoine-langlois.net/datahearth/config-mapper/internal/providers"
)

func Load(cfg *cmConfig.Definition, providers ...cmProviders.Provider) error {
	return nil
}
