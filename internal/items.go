package mapper

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/viper"
)

type Items struct {
	locations  []ItemLocation
	storage    string
	repository RepositoryActions
	indexer    Indexer
}

type ItemsActions interface {
	Action(action string)
	AddItems(items []ItemLocation)
	CleanUp(removedLines []string) error
}

func NewItemsActions(items []ItemLocation, storage string, repository RepositoryActions, indexer Indexer) ItemsActions {
	if items == nil {
		items = []ItemLocation{}
	}

	return &Items{
		locations:  items,
		storage:    storage,
		repository: repository,
		indexer:    indexer,
	}
}

func (e *Items) Action(a string) {
	color.Blue("# %s", a)
	newLines := []string{}

	for i, l := range e.locations {
		var src, dst string
		storagePath, systemPath, err := configPaths(l, e.storage)
		if err != nil {
			PrintError("[%d] failed to resolve item paths \"%v\": %v", i, l, err)
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
			PrintError("[%d] failed to create directory architecture for destination path \"%s\": %v", i, path.Dir(dst), err)
			continue
		}

		s, err := os.Stat(src)
		if err != nil {
			PrintError("[%d] failed to check if source path is a folder \"%s\": %v", i, src, err)
			continue
		}

		if s.IsDir() {
			dstPerms := fs.FileMode(0755)
			s, err := os.Stat(dst)
			if err != nil {
				if !os.IsNotExist(err) {
					PrintError("[%d] failed to check if destination folder \"%s\" exists: %v", i, dst, err)
					continue
				}
			} else {
				dstPerms = s.Mode()
			}

			if err := os.Mkdir(dst, dstPerms); err != nil {
				if !os.IsExist(err) {
					PrintError("[%d] failed to create destination folder \"%s\": %v", i, dst, err)
					continue
				}
			}
			if err := copyFolder(src, dst); err != nil {
				PrintError("[%d] failed to %s folder from \"%s\" to \"%s\": %v", i, a, src, dst, err)
				continue
			}
		} else {
			if err := copyFile(src, dst); err != nil {
				PrintError("[%d] failed to %s file from \"%s\" to \"%s\": %v", i, a, src, dst, err)
				continue
			}
		}

		if a == "save" {
			p, err := absolutePath(e.storage)
			if err != nil {
				PrintError("[%d] failed resolve absolute path from configuration storage: %v", i, err)
				continue
			}
			newLines = append(newLines, strings.ReplaceAll(dst, p+"/", ""))
		}

		color.Green("[%d] %s copied", i, src)
	}

	if a == "save" && !viper.GetBool("disable-index-update") {
		if err := e.indexer.Write(newLines); err != nil {
			PrintError(err.Error())
			os.Exit(1)
		}
	}
}

func (e *Items) AddItems(items []ItemLocation) {
	e.locations = append(e.locations, items...)
}

func (e *Items) CleanUp(removedLines []string) error {
	for _, l := range removedLines {
		path, err := absolutePath(fmt.Sprintf("%s/%s", e.storage, l))
		if err != nil {
			return err
		}

		if err := os.RemoveAll(path); err != nil {
			return fmt.Errorf("failed to remove item %s: %v", l, err)
		}
	}

	return nil
}
