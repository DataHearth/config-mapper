package mapper

import (
	"fmt"
	"io/fs"
	"os"
	"strings"
)

type Index struct {
	lines        []string
	path         string
	perms        fs.FileMode
	repoPath     string
	removedLines []string
}

type Indexer interface {
	Write(newLines []string) error
	filter(configLines []string) map[string]bool
	RemovedLines() []string
	Lines() []string
}

func NewIndexer(repoPath string) (Indexer, error) {
	perms := fs.FileMode(0755)
	indexPath, err := absolutePath(fmt.Sprintf("%s/.index", repoPath))
	if err != nil {
		return nil, err
	}

	var l []string
	s, err := os.Stat(indexPath)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}

		l = []string{}
	} else {
		perms = s.Mode()

		b, err := os.ReadFile(indexPath)
		if err != nil {
			return nil, err
		}

		l = strings.Split(string(b), "\n")
	}

	return &Index{
		lines:        l,
		path:         indexPath,
		perms:        perms,
		repoPath:     repoPath,
		removedLines: []string{},
	}, nil
}

func (i *Index) RemovedLines() []string {
	return i.removedLines
}

func (i *Index) Lines() []string {
	return i.lines
}

// filter removes lines that are no more used in configuration from the index
func (i *Index) filter(newLines []string) map[string]bool {
	removedLines := []string{}
	foundLines := map[string]bool{}
	for _, nl := range newLines {
		foundLines[nl] = true
	}
	for _, ml := range i.lines {
		if _, ok := foundLines[ml]; !ok {
			removedLines = append(removedLines, ml)
		}
	}

	i.removedLines = removedLines
	return foundLines
}

// Write add lines stored in memory the .index file
func (i *Index) Write(newLines []string) error {
	lines := i.filter(newLines)

	var data []byte
	index := 0
	linesNumber := len(lines)
	for path := range lines {
		if index+1 == linesNumber {
			data = append(data, []byte(fmt.Sprint(path))...)
		} else {
			data = append(data, []byte(fmt.Sprintln(path))...)
		}

		index += 1
	}

	os.WriteFile(i.path, data, i.perms)

	return nil
}
