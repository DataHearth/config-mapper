package mapper

import (
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/pterm/pterm"
)

var (
	ErrCopy          = errors.New("failed to copy some files")
	ErrUnsupportedOS = errors.New("unsupported OS. Please, contact the maintainer")
)

func LoadFiles(files []ItemLocation, location string) error {
	pterm.DefaultSection.Println("Save files into saved location")
	haveErr := false
	p, _ := pterm.DefaultProgressbar.WithTotal(len(files)).Start()

	for _, f := range files {
		src, dst, err := configPaths(f, location)
		if err != nil {
			if err == ErrUnsupportedOS {
				return err
			}
			pterm.Error.Println(fmt.Sprintf("failed to destination resolve path \"%s\": %v", f.Linux, err))
			haveErr = true
			continue
		}

		if err := os.MkdirAll(path.Dir(dst), 0755); err != nil {
			pterm.Error.Printfln(fmt.Sprintf("failed to create directory architecture for destination path \"%s\": %v", dst, err))
			haveErr = true
			continue
		}

		p.UpdateTitle(fmt.Sprintf("copying %s", src))

		if err := copyFile(src, dst); err != nil {
			pterm.Error.Println(fmt.Sprintf("failed to load file from \"%s\" to \"%s\": %v", src, dst, err))
			haveErr = true
			continue
		}

		pterm.Success.Println(fmt.Sprintf("%s copied", src))
		p.Increment()
	}

	p.Stop()
	if haveErr {
		return ErrCopy
	}

	return nil
}

func SaveFiles(files []ItemLocation, location string) error {
	haveErr := false
	pterm.DefaultSection.Println("Save files into saved location")
	p, _ := pterm.DefaultProgressbar.WithTotal(len(files)).Start()

	for _, f := range files {
		dst, src, err := configPaths(f, location)
		if err != nil {
			if err == ErrUnsupportedOS {
				return err
			}
			pterm.Error.Println(fmt.Sprintf("failed to destination resolve path \"%s\": %v", f.Linux, err))
			haveErr = true
			continue
		}

		p.UpdateTitle(fmt.Sprintf("copying \"%s\"", src))

		if err := os.MkdirAll(path.Dir(dst), 0755); err != nil {
			pterm.Error.Printfln(fmt.Sprintf("failed to create directory architecture for destination path \"%s\": %v", dst, err))
			haveErr = true
			continue
		}

		if err := copyFile(src, dst); err != nil {
			pterm.Error.Println(fmt.Sprintf("failed to load file from \"%s\" to \"%s\": %v", src, dst, err))
			haveErr = true
			continue
		}

		pterm.Success.Println(fmt.Sprintf("\"%s\" copied", src))
		p.Increment()
	}

	p.Stop()
	if haveErr {
		return ErrCopy
	}

	return nil
}
