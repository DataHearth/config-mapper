package mapper

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"runtime"

	"github.com/pterm/pterm"
)

var (
	ErrCopy          = errors.New("failed to copy some files")
	ErrUnsupportedOS = errors.New("unsupported OS. Please, contact the maintainer")
)

func LoadFiles(files []ItemLocation, location string) error {
	haveErr := false
	p, _ := pterm.DefaultProgressbar.WithTotal(len(files)).WithTitle("Loading files onto your system").Start()

	for _, f := range files {
		var src, dst string
		var err error

		switch runtime.GOOS {
		case "linux":
			src, dst, err = getPaths(f.Linux, location)
			if err != nil {
				pterm.Error.Println(fmt.Sprintf("failed to destination resolve path \"%s\": %v", f.Linux, err))
				haveErr = true
				continue
			}
		case "darwin":
			src, dst, err = getPaths(f.Darwin, location)
			if err != nil {
				pterm.Error.Println(fmt.Sprintf("failed to destination resolve path \"%s\": %v", f.Darwin, err))
				haveErr = true
				continue
			}
		default:
			return ErrUnsupportedOS
		}

		if err := os.MkdirAll(path.Dir(dst), 0755); err != nil {
			pterm.Error.Printfln(fmt.Sprintf("failed to create directory architecture for destination path \"%s\": %v", dst, err))
			haveErr = true
			continue
		}

		p.UpdateTitle(fmt.Sprintf("copying %s", src))

		if err := copy(src, dst); err != nil {
			pterm.Error.Println(fmt.Sprintf("failed to load file from \"%s\" to \"%s\": %v", src, dst, err))
			haveErr = true
			continue
		}

		pterm.Success.Println(fmt.Sprintf("%s copied", src))
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
		var src, dst string
		var err error

		p.UpdateTitle(fmt.Sprintf("copying %s", src))

		switch runtime.GOOS {
		case "linux":
			dst, src, err = getPaths(f.Linux, location)
			if err != nil {
				pterm.Error.Println(fmt.Sprintf("failed to destination resolve path \"%s\": %v", f.Linux, err))
				haveErr = true
				continue
			}
		case "darwin":
			dst, src, err = getPaths(f.Darwin, location)
			if err != nil {
				pterm.Error.Println(fmt.Sprintf("failed to destination resolve path \"%s\": %v", f.Darwin, err))
				haveErr = true
				continue
			}
		default:
			return ErrUnsupportedOS
		}

		if err := os.MkdirAll(path.Dir(dst), 0755); err != nil {
			pterm.Error.Printfln(fmt.Sprintf("failed to create directory architecture for destination path \"%s\": %v", dst, err))
			haveErr = true
			continue
		}

		if err := copy(src, dst); err != nil {
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

func copy(src, dst string) error {
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
