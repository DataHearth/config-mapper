package mapper

import (
	"fmt"
	"io"
	"os"
	"path"
	"runtime"
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

func copyFile(src, dst string) error {
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

func configPaths(f ItemLocation, location string) (string, string, error) {
	var src, dst string
	var err error

	switch runtime.GOOS {
	case "linux":
		src, dst, err = getPaths(f.Linux, location)
		if err != nil {
			return "", "", err
		}
	case "darwin":
		src, dst, err = getPaths(f.Darwin, location)
		if err != nil {
			return "", "", err
		}
	default:
		return "", "", ErrUnsupportedOS
	}

	return src, dst, nil
}

func copyFolder(src, dst string) error {
	items, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, i := range items {
		itemName := i.Name()
		srcItem := fmt.Sprintf("%s/%s", src, itemName)
		dstItem := fmt.Sprintf("%s/%s", dst, itemName)

		if i.IsDir() {
			if err := os.Mkdir(dstItem, i.Type().Perm()); err != nil {
				return err
			}
			if err := copyFolder(srcItem, dstItem); err != nil {
				return err
			}
		}

		if err := copyFile(srcItem, dstItem); err != nil {
			return err
		}
	}

	return nil
}
