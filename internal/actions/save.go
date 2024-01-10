package actions

import (
	"io"
	"os"
	"path"
	"strings"

	cmConfig "gitea.antoine-langlois.net/datahearth/config-mapper/internal/config"
	cmProviders "gitea.antoine-langlois.net/datahearth/config-mapper/internal/providers"
)

func Save(logger, cfg *cmConfig.Definition, providers ...cmProviders.Provider) error {
	dstDir := path.Join(cfg.Path, cfg.UUID.String())
	for _, i := range cfg.Items {
		paths := strings.Split(i, ":")
		in, out := paths[0], path.Join(dstDir, paths[1])

		s, err := os.Stat(in)
		if err != nil {
			return err
		}

		if s.IsDir() {
			if err := copyDir(in, out); err != nil {
				return err
			}

			continue
		}

		if err := copyFile(in, out); err != nil {
			return err
		}
	}

	return nil
}

func copyDir(in string, out string) error {
	entries, err := os.ReadDir(in)
	if err != nil {
		return err
	}

	for _, e := range entries {
		if e.IsDir() {
			return copyDir(path.Join(in, e.Name()), path.Join(out, e.Name()))
		}

		if err := copyFile(path.Join(in, e.Name()), path.Join(out, e.Name())); err != nil {
			return err
		}
	}

	return nil
}

func copyFile(in string, out string) error {
	fIn, err := os.Open(in)
	if err != nil {
		return err
	}
	defer fIn.Close()

	sIn, err := fIn.Stat()
	if err != nil {
		return err
	}

	fOut, err := os.OpenFile(out, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, sIn.Mode())
	if err != nil {
		return err
	}
	defer fOut.Close()

	if _, err := io.Copy(fOut, fIn); err != nil {
		return err
	}

	return nil
}
