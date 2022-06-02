package mapper

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"strings"

	"github.com/datahearth/config-mapper/internal/configuration"
	"github.com/datahearth/config-mapper/internal/git"
	"github.com/datahearth/config-mapper/internal/misc"
	"github.com/fatih/color"
	"github.com/spf13/viper"
)

type Items struct {
	locations  []configuration.OSLocation
	storage    string
	repository git.RepositoryActions
	indexer    Indexer
}

type ItemsActions interface {
	Action(action string)
	AddItems(items []configuration.OSLocation)
	CleanUp(removedLines []string) error
}

func NewItemsActions(items []configuration.OSLocation, storage string, repository git.RepositoryActions, indexer Indexer) ItemsActions {
	if items == nil {
		items = []configuration.OSLocation{}
	}

	return &Items{
		locations:  items,
		storage:    storage,
		repository: repository,
		indexer:    indexer,
	}
}

func (e *Items) Action(action string) {
	color.Blue("# %s", action)
	newLines := []string{}

	for i, l := range e.locations {
		var src string
		storagePath, systemPath, err := misc.ConfigPaths(l, e.storage)
		if err != nil {
			PrintError("[%d] failed to resolve item paths \"%v\": %v", i, l, err)
			continue
		}

		if action == "save" {
			src = systemPath

			if newItem := e.saveItem(systemPath, storagePath, i); newItem != "" {
				newLines = append(newLines, newItem)
			} else {
				continue
			}
		} else {
			src = storagePath
			e.loadItem(storagePath, systemPath, i)
		}

		color.Green("[%d] %s copied", i, src)
	}

	if action == "save" && !viper.GetBool("disable-index-update") {
		if err := e.indexer.Write(newLines); err != nil {
			PrintError(err.Error())
			os.Exit(1)
		}
	}
}

func (e *Items) saveItem(src, dst string, index int) string {
	if err := os.MkdirAll(path.Dir(dst), 0755); err != nil {
		PrintError("[%d] failed to create directory architecture for destination path \"%s\": %v", index, path.Dir(dst), err)
		return ""
	}

	s, err := os.Stat(src)
	if err != nil {
		PrintError("[%d] failed to check if source path is a folder \"%s\": %v", index, src, err)
		return ""
	}

	if s.IsDir() {
		dstPerms := fs.FileMode(0755)
		s, err := os.Stat(dst)
		if err != nil {
			if !os.IsNotExist(err) {
				PrintError("[%d] failed to check if destination folder \"%s\" exists: %v", index, dst, err)
				return ""
			}
		} else {
			dstPerms = s.Mode()
		}

		if err := os.Mkdir(dst, dstPerms); err != nil {
			if !os.IsExist(err) {
				PrintError("[%d] failed to create destination folder \"%s\": %v", index, dst, err)
				return ""
			}
		}
		if err := misc.CopyFolder(src, dst, true); err != nil {
			PrintError("[%d] failed to save folder from \"%s\" to \"%s\": %v", index, src, dst, err)
			return ""
		}
	} else {
		if err := misc.CopyFile(src, dst); err != nil {
			PrintError("[%d] failed to save file from \"%s\" to \"%s\": %v", index, src, dst, err)
			return ""
		}
	}

	p, err := misc.AbsolutePath(e.storage)
	if err != nil {
		PrintError("[%d] failed resolve absolute path from configuration storage: %v", index, err)
		return ""
	}

	return strings.ReplaceAll(dst, p+"/", "")
}

func (e *Items) loadItem(src, dst string, index int) {
	if err := os.MkdirAll(path.Dir(dst), 0755); err != nil {
		PrintError("[%d] failed to create directory architecture for destination path \"%s\": %v", index, path.Dir(dst), err)
		return
	}

	s, err := os.Stat(src)
	if err != nil {
		PrintError("[%d] failed to check if source path is a folder \"%s\": %v", index, src, err)
		return
	}

	if s.IsDir() {
		dstPerms := fs.FileMode(0755)
		s, err := os.Stat(dst)
		if err != nil {
			if !os.IsNotExist(err) {
				PrintError("[%d] failed to check if destination folder \"%s\" exists: %v", index, dst, err)
				return
			}
		} else {
			dstPerms = s.Mode()
		}

		if err := os.Mkdir(dst, dstPerms); err != nil {
			if !os.IsExist(err) {
				PrintError("[%d] failed to create destination folder \"%s\": %v", index, dst, err)
				return
			}
		}
		if err := misc.CopyFolder(src, dst, false); err != nil {
			PrintError("[%d] failed to load folder from \"%s\" to \"%s\": %v", index, src, dst, err)
			return
		}
	} else {
		if err := misc.CopyFile(src, dst); err != nil {
			PrintError("[%d] failed to load file from \"%s\" to \"%s\": %v", index, src, dst, err)
			return
		}
	}
}

func (e *Items) AddItems(items []configuration.OSLocation) {
	e.locations = append(e.locations, items...)
}

func (e *Items) CleanUp(removedLines []string) error {
	for _, l := range removedLines {
		path, err := misc.AbsolutePath(fmt.Sprintf("%s/%s", e.storage, l))
		if err != nil {
			return err
		}

		if err := os.RemoveAll(path); err != nil {
			return fmt.Errorf("failed to remove item %s: %v", l, err)
		}
	}

	return nil
}

func PrintError(err string, values ...interface{}) {
	color.Error.Write([]byte(color.RedString(err+"\n", values...)))
}
