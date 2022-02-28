package mapper

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/pterm/pterm"
)

var (
	ErrCopy          = errors.New("failed to copy some files")
	ErrUnsupportedOS = errors.New("unsupported OS. Please, contact the maintainer")
)

type Items struct {
	locations  []ItemLocation
	storage    string
	progessBar *pterm.ProgressbarPrinter
	runErr     bool
}

type ItemsActions interface {
	Action(action string) error
	AddItems(items []ItemLocation)
	copy(src, dst string) error
}

func NewElement(l []ItemLocation, storage string) ItemsActions {
	return &Items{
		l,
		storage,
		pterm.DefaultProgressbar.WithTotal(len(l)),
		false,
	}
}

func (e *Items) Action(a string) error {
	pterm.DefaultSection.Println(fmt.Sprintf("%s items", strings.Title(a)))

	e.progessBar.Start()

	for _, f := range e.locations {
		var src, dst string
		storagePath, systemPath, err := configPaths(f, e.storage)
		if err != nil {
			pterm.Error.Println(fmt.Sprintf("failed to resolve item paths \"%v\": %v", f, err))
			e.runErr = true
			continue
		}
		if a == "save" {
			src = systemPath
			dst = storagePath
		} else {
			src = storagePath
			dst = systemPath
		}

		if err := os.MkdirAll(path.Dir(dst), 0755); err != nil {
			pterm.Error.Printfln(fmt.Sprintf("failed to create directory architecture for destination path \"%s\": %v", path.Dir(dst), err))
			e.runErr = true
			continue
		}

		s, err := os.Stat(src)
		if err != nil {
			pterm.Error.Println(fmt.Sprintf("failed to check if source path is a folder \"%s\": %v", src, err))
			e.runErr = true
			continue
		}

		e.progessBar.UpdateTitle(fmt.Sprintf("copying %s", src))

		if s.IsDir() {
			if err := os.Mkdir(dst, 0755); err != nil {
				if !os.IsExist(err) {
					pterm.Error.Println(fmt.Sprintf("failed to create destination folder \"%s\": %v", dst, err))
					e.runErr = true
					continue
				}
			}
			if err := copyFolder(src, dst); err != nil {
				pterm.Error.Println(fmt.Sprintf("failed to %s folder from \"%s\" to \"%s\": %v", a, src, dst, err))
				e.runErr = true
				continue
			}
		} else {
			if err := copyFile(src, dst); err != nil {
				pterm.Error.Println(fmt.Sprintf("failed to %s file from \"%s\" to \"%s\": %v", a, src, dst, err))
				e.runErr = true
				continue
			}
		}

		pterm.Success.Println(fmt.Sprintf("%s copied", src))
		e.progessBar.Increment()
	}

	e.progessBar.Stop()
	if e.runErr {
		e.runErr = false
		return ErrCopy
	}

	return nil
}

func (e *Items) AddItems(items []ItemLocation) {
	e.locations = append(e.locations, items...)
}

func (e *Items) copy(src, dst string) error {
	s, err := os.Stat(src)
	if err != nil {
		return err
	}

	e.progessBar.UpdateTitle(fmt.Sprintf("copying %s", src))

	if s.IsDir() {
		if err := copyFolder(src, dst); err != nil {
			return err
		}
	} else {
		if err := copyFile(src, dst); err != nil {
			return err
		}
	}

	pterm.Success.Println(fmt.Sprintf("\"%s\" copied", src))
	e.progessBar.Increment()

	return nil
}
