package internal

import (
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

// ResolvePath resolves the path using environment variables and "~"
func ResolvePath(path string) (string, error) {
	path = os.ExpandEnv(path)
	if strings.Contains(path, "~") {
		usr, err := user.Current()
		if err != nil {
			return "", err
		}

		path = filepath.Join(usr.HomeDir, strings.Replace(path, "~", "", 1))
	}

	return path, nil
}
