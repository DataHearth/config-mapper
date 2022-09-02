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

// Action performs a "save" or "load" action on all given items.
//
// Any error is printed to STDERR and item is skipped.
//
// If the performed action is "save", it'll also write the `.index` file with all new items.
func (e *Items) Action(action string) {
	fmt.Printf("# %s files and folders\n", action)
	newLines := []string{}

	for i, l := range e.locations {
		var src string
		storagePath, systemPath, err := misc.ConfigPaths(l, e.storage)
		if err != nil {
			PrintError("❌ failed to resolve item paths \"%v\": %v", i, l, err)
			continue
		}
		if storagePath == "" && systemPath == "" {
			fmt.Println("⛔ Skipping")
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

		fmt.Printf("✔️ %s\n", src)
	}

	if action == "save" && !viper.GetBool("disable-index-update") {
		if err := e.indexer.Write(newLines); err != nil {
			PrintError(err.Error())
			os.Exit(1)
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
		PrintError("❌ failed to create directory architecture for destination path \"%s\": %v", index, path.Dir(dst), err)
		return ""
	}

	s, err := os.Stat(src)
	if err != nil {
		PrintError("❌ failed to check if source path is a folder \"%s\": %v", index, src, err)
		return ""
	}

	if s.IsDir() {
		dstPerms := fs.FileMode(0755)
		s, err := os.Stat(dst)
		if err != nil {
			if !os.IsNotExist(err) {
				PrintError("❌ failed to check if destination folder \"%s\" exists: %v", index, dst, err)
				return ""
			}
		} else {
			dstPerms = s.Mode()
		}

		// remove the destination if it exists. It cleans up the saved location from unused files
		if err := os.RemoveAll(dst); err != nil {
			PrintError("❌ failed to truncate destination folder \"%s\": %v", index, dst, err)
		}

		if err := os.Mkdir(dst, dstPerms); err != nil {
			if !os.IsExist(err) {
				PrintError("❌ failed to create destination folder \"%s\": %v", index, dst, err)
				return ""
			}
		}
		if err := misc.CopyFolder(src, dst, true); err != nil {
			PrintError("❌ failed to save folder from \"%s\" to \"%s\": %v", index, src, dst, err)
			return ""
		}
	} else {
		if err := misc.CopyFile(src, dst); err != nil {
			PrintError("❌ failed to save file from \"%s\" to \"%s\": %v", index, src, dst, err)
			return ""
		}
	}

	p, err := misc.AbsolutePath(e.storage)
	if err != nil {
		PrintError("❌ failed resolve absolute path from configuration storage: %v", index, err)
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
		PrintError("❌ failed to create directory architecture for destination path \"%s\": %v", index, path.Dir(dst), err)
		return
	}

	s, err := os.Stat(src)
	if err != nil {
		PrintError("❌ failed to check if source path is a folder \"%s\": %v", index, src, err)
		return
	}

	if s.IsDir() {
		dstPerms := fs.FileMode(0755)
		s, err := os.Stat(dst)
		if err != nil {
			if !os.IsNotExist(err) {
				PrintError("❌ failed to check if destination folder \"%s\" exists: %v", index, dst, err)
				return
			}
		} else {
			dstPerms = s.Mode()
		}

		if err := os.Mkdir(dst, dstPerms); err != nil {
			if !os.IsExist(err) {
				PrintError("❌ failed to create destination folder \"%s\": %v", index, dst, err)
				return
			}
		}
		if err := misc.CopyFolder(src, dst, false); err != nil {
			PrintError("❌ failed to load folder from \"%s\" to \"%s\": %v", index, src, dst, err)
			return
		}
	} else {
		if err := misc.CopyFile(src, dst); err != nil {
			PrintError("❌ failed to load file from \"%s\" to \"%s\": %v", index, src, dst, err)
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
