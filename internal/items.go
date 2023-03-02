package mapper

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"strings"

	"gitea.antoine-langlois.net/datahearth/config-mapper/internal/configuration"
	"gitea.antoine-langlois.net/datahearth/config-mapper/internal/git"
	"gitea.antoine-langlois.net/datahearth/config-mapper/internal/misc"
	"github.com/charmbracelet/log"
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

// Action performs a "save" or "load" action on all given items.
//
// Any error is printed to STDERR and item is skipped.
//
// If the performed action is "save", it'll also write the `.index` file with all new items.
func (e *Items) Action(action string) {
	log.Info("performing action", "action", action)
	newLines := []string{}

	for i, l := range e.locations {
		storagePath, systemPath, err := misc.ConfigPaths(l, e.storage)
		if err != nil {
			log.Error("failed to resolve item paths", "item", i, "location", l, "err", err)
			continue
		}
		if storagePath == "" && systemPath == "" {
			log.Info("item is empty", "item", i, "location", l)
			continue
		}

		if action == "save" {
			if newItem := e.saveItem(systemPath, storagePath, i); newItem != "" {
				newLines = append(newLines, newItem)
			} else {
				continue
			}
		} else {
			e.loadItem(storagePath, systemPath, i)
		}

		log.Info("item processed", "action", action, "item", i, "location", l)
	}

	if action == "save" && !viper.GetBool("disable-index-update") {
		if err := e.indexer.Write(newLines); err != nil {
			log.Fatal(err)
		}
	}
}

// saveItem saves a given item inside the configured saved location.
//
// If an error is given during the process, the function returns an empty string
// (meaning the item hasn't been saved) and prints the error in STDERR.
//
// Else, returns the relative item location from the saved location to write the index
// (E.g: /home/user/.config => .config)
func (e *Items) saveItem(src, dst string, index int) string {
	if err := os.MkdirAll(path.Dir(dst), 0755); err != nil {
		log.Error("failed to create directory architecture for destination path", "path", path.Dir(dst), "err", err)
		return ""
	}

	s, err := os.Stat(src)
	if err != nil {
		log.Error("failed to check if source path is a folder", "path", src, "err", err)
		return ""
	}

	if s.IsDir() {
		dstPerms := fs.FileMode(0755)
		s, err := os.Stat(dst)
		if err != nil {
			if !os.IsNotExist(err) {
				log.Error("failed to check if destination folder exists", "path", dst, "err", err)
				return ""
			}
		} else {
			dstPerms = s.Mode()
		}

		// remove the destination if it exists. It cleans up the saved location from unused files
		if err := os.RemoveAll(dst); err != nil {
			log.Error("failed to truncate destination folder", "path", dst, "err", err)
		}

		if err := os.Mkdir(dst, dstPerms); err != nil {
			if !os.IsExist(err) {
				log.Error("failed to create destination folder", "path", dst, "err", err)
				return ""
			}
		}
		if err := misc.CopyFolder(src, dst, true); err != nil {
			log.Error("failed to save folder from source to destination", "source", src, "destination", dst, "err", err)
			return ""
		}
	} else {
		if err := misc.CopyFile(src, dst); err != nil {
			log.Error("failed to save file from source to destination", "source", src, "destination", dst, "err", err)
			return ""
		}
	}

	p, err := misc.AbsolutePath(e.storage)
	if err != nil {
		log.Error("failed resolve absolute path from configuration storage", "err", err)
		return ""
	}

	return strings.ReplaceAll(dst, p+"/", "")
}

// loadItem loads a given item onto the system.
//
// If an error is given during the process, the function returns an empty string
// (meaning the item hasn't been saved) and prints the error in STDERR.
func (e *Items) loadItem(src, dst string, index int) {
	if err := os.MkdirAll(path.Dir(dst), 0755); err != nil {
		log.Error("failed to create directory architecture for destination path", "path", path.Dir(dst), "err", err)
		return
	}

	s, err := os.Stat(src)
	if err != nil {
		log.Error("failed to check if source path is a folder", "path", src, "err", err)
		return
	}

	if s.IsDir() {
		dstPerms := fs.FileMode(0755)
		s, err := os.Stat(dst)
		if err != nil {
			if !os.IsNotExist(err) {
				log.Error("failed to check if destination folder exists", "path", dst, "err", err)
				return
			}
		} else {
			dstPerms = s.Mode()
		}

		if err := os.Mkdir(dst, dstPerms); err != nil {
			if !os.IsExist(err) {
				log.Error("failed to create destination folder", "path", dst, "err", err)
				return
			}
		}
		if err := misc.CopyFolder(src, dst, false); err != nil {
			log.Error("failed to load folder from source to destination", "source", src, "destination", dst, "err", err)
			return
		}
	} else {
		if err := misc.CopyFile(src, dst); err != nil {
			log.Error("failed to load file from source to destination", "source", src, "destination", dst, "err", err)
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
