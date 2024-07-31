package misc

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"runtime"
	"strings"

	"github.com/datahearth/config-mapper/internal/configuration"
)

func AbsolutePath(p string) (string, error) {
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
				return "", errors.New("found invalid environment variable in path")

			}
			pathPart = env
		}

		finalPath += fmt.Sprintf("/%s", pathPart)
	}

	return path.Clean(finalPath), nil
}

func getPaths(p string, l string) (string, string, error) {
	paths := strings.Split(p, ":")

	if len(paths) < 2 {
		return "", "", errors.New("path incorrectly formatted. It requires \"source:destination\"")
	}

	src, err := AbsolutePath(strings.Replace(paths[0], "$LOCATION", l, 1))
	if err != nil {
		return "", "", err
	}

	dst, err := AbsolutePath(paths[1])
	if err != nil {
		return "", "", err
	}

	return src, dst, nil
}

func CopyFile(src, dst string) error {
	s, err := os.Stat(src)
	if err != nil {
		return err
	}

	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		return err
	}

	if err := os.Chmod(dst, s.Mode()); err != nil {
		return err
	}

	return nil
}

func ConfigPaths(os configuration.OSLocation, location string) (string, string, error) {
	var src, dst string
	var err error

	switch runtime.GOOS {
	case "linux":
		if os.Linux == "" {
			return "", "", nil
		}
		src, dst, err = getPaths(os.Linux, location)
		if err != nil {
			return "", "", err
		}
	case "darwin":
		if os.Darwin == "" {
			return "", "", nil
		}
		src, dst, err = getPaths(os.Darwin, location)
		if err != nil {
			return "", "", err
		}
	default:
		return "", "", errors.New("unsupported OS. Please, contact the maintainer")
	}

	return src, dst, nil
}

var ignored map[string]bool

func CopyFolder(src, dst string, checkIgnore bool) error {
	items, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	if checkIgnore {
		f, err := os.ReadFile(fmt.Sprintf("%s/.ignore", src))
		if err != nil && !errors.Is(err, io.EOF) {
			if !errors.Is(err, os.ErrNotExist) {
				return err
			}
		}

		ignored = map[string]bool{}
		for _, l := range strings.Split(string(f), "\n") {
			if l != "" && !strings.Contains(l, "#") {
				ignored[fmt.Sprintf("%s/%s", src, l)] = true
			}
		}
	}

	for _, i := range items {
		itemName := i.Name()
		srcItem := fmt.Sprintf("%s/%s", src, itemName)
		// do not copy item if it's present in .ignore file
		if ignored != nil {
			if _, ok := ignored[srcItem]; ok {
				continue
			}
		}

		dstItem := fmt.Sprintf("%s/%s", dst, itemName)

		if i.IsDir() {
			info, err := i.Info()
			if err != nil {
				return err
			}

			if err := os.MkdirAll(dstItem, info.Mode()); err != nil {
				return err
			}
			if err := CopyFolder(srcItem, dstItem, false); err != nil {
				return err
			}

			continue
		}

		if err := CopyFile(srcItem, dstItem); err != nil {
			return err
		}
	}

	return nil
}
