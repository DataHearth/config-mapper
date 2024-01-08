package git

import "gitea.antoine-langlois.net/datahearth/config-mapper/internal/providers"

type Git struct {
	rootDir   string
	remoteDir string
}

func New(rootDir string, remoteDir string) providers.Provider {
	return Git{
		rootDir:   rootDir,
		remoteDir: remoteDir,
	}
}

func (g Git) Upload() error {
	return nil
}

func (g Git) Download() error {
	return nil
}
