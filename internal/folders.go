package mapper

import (
	"errors"
	"fmt"
	"os"

	"github.com/pterm/pterm"
)

var ErrFolderCopy = errors.New("failed to copy some folders")

func LoadFolders(folders []ItemLocation, location string) error {
	haveErr := false
	pterm.DefaultSection.Println("Load folders into saved location")
	p, _ := pterm.DefaultProgressbar.WithTotal(len(folders)).Start()

	for _, f := range folders {
		src, dst, err := configPaths(f, location)
		if err != nil {
			if err == ErrUnsupportedOS {
				return err
			}
			pterm.Error.Println(fmt.Sprintf("failed to destination resolve path \"%s\": %v", f.Linux, err))
			haveErr = true
			continue
		}

		s, err := os.Stat(src)
		if err != nil {
			pterm.Error.Println(fmt.Sprintf("failed to check if source path is a folder \"%s\": %v", src, err))
			haveErr = true
			continue
		}
		if !s.IsDir() {
			pterm.Error.Println(fmt.Sprintf("source path is a file \"%s\"", src))
			haveErr = true
			continue
		}

		p.UpdateTitle(fmt.Sprintf("copying folder \"%s\"", src))

		if err := os.MkdirAll(dst, 0755); err != nil {
			pterm.Error.Printfln(fmt.Sprintf("failed to create directory architecture for destination path \"%s\": %v", dst, err))
			haveErr = true
			continue
		}

		if err := copyFolder(src, dst); err != nil {
			pterm.Error.Println(fmt.Sprintf("failed to load folder from \"%s\" to \"%s\": %v", src, dst, err))
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

func SaveFolders(folders []ItemLocation, location string) error {
	haveErr := false
	pterm.DefaultSection.Println("Save folders into saved location")
	p, _ := pterm.DefaultProgressbar.WithTotal(len(folders)).Start()

	for _, f := range folders {
		dst, src, err := configPaths(f, location)
		if err != nil {
			if err == ErrUnsupportedOS {
				return err
			}
			pterm.Error.Println(fmt.Sprintf("failed to destination resolve path \"%s\": %v", f.Linux, err))
			haveErr = true
			continue
		}

		s, err := os.Stat(src)
		if err != nil {
			pterm.Error.Println(fmt.Sprintf("failed to check if source path is a folder \"%s\": %v", src, err))
			haveErr = true
			continue
		}
		if !s.IsDir() {
			pterm.Error.Println(fmt.Sprintf("source path is a file \"%s\"", src))
			haveErr = true
			continue
		}

		p.UpdateTitle(fmt.Sprintf("copying folder \"%s\"", src))

		if err := os.MkdirAll(dst, 0755); err != nil {
			pterm.Error.Printfln(fmt.Sprintf("failed to create directory architecture for destination path \"%s\": %v", dst, err))
			haveErr = true
			continue
		}

		if err := copyFolder(src, dst); err != nil {
			pterm.Error.Println(fmt.Sprintf("failed to save folder from \"%s\" to \"%s\": %v", src, dst, err))
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

func copyFolder(src, dst string) error {
	var haveErr bool

	items, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, i := range items {
		itemName := i.Name()
		srcItem := fmt.Sprintf("%s/%s", src, itemName)
		dstItem := fmt.Sprintf("%s/%s", dst, itemName)

		if i.IsDir() {
			os.Mkdir(dstItem, i.Type().Perm())
			copyFolder(srcItem, dstItem)
			continue
		}

		if err := copyFile(srcItem, dstItem); err != nil {
			haveErr = true
			continue
		}
	}

	if haveErr {
		return ErrFolderCopy
	}
	return nil
}
