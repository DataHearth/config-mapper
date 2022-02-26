package mapper

import (
	"fmt"
	"os"
	"path"
	"strings"
)

func absolutePath(p string) (string, error) {
	finalPath := p
	if strings.Contains(finalPath, "~") {
		h, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}

		finalPath = strings.Replace(p, "~", h, 1)
	}

	splitted := strings.Split(finalPath, "/")
	finalPath = ""
	for _, s := range splitted {
		pathPart := s
		if strings.Contains(s, "$") {
			env := os.Getenv(s)
			if env == "" {
				return "", ErrInvalidEnv
			}
			pathPart = env
		}

		finalPath += fmt.Sprintf("/%s", pathPart)
	}

	return path.Clean(finalPath), nil
}

func getPaths(p string, l string) (string, string, error) {
	paths := strings.Split(p, ":")

	src, err := absolutePath(strings.Replace(paths[0], "$LOCATION", l, 1))
	if err != nil {
		return "", "", err
	}

	dst, err := absolutePath(paths[1])
	if err != nil {
		return "", "", err
	}

	return src, dst, nil
}
