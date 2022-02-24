package mapper

import (
	"errors"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/spf13/viper"
)

var (
	ErrDirIsFile      = errors.New("path is a file")
	ErrEmptyGitConfig = errors.New("empty git configuration")
)

func OpenGitRepo() (*git.Repository, error) {
	configFolder := viper.GetString("storage.location")

	s, err := os.Stat(configFolder)
	if err != nil {
		if os.IsNotExist(err) {
			gitConfig := viper.GetStringMapString("storage.git")
			if gitConfig == nil {
				return nil, ErrEmptyGitConfig
			}

			repo, err := git.PlainClone(viper.GetString("storage.location"), false, &git.CloneOptions{
				URL:      gitConfig["repository"],
				Progress: os.Stdout,
				Auth: &http.BasicAuth{
					Username: gitConfig["username"],
					Password: gitConfig["password"],
				},
			})
			if err != nil {
				return nil, err
			}

			return repo, nil
		}

		return nil, err
	}

	if s.IsDir() {
		return nil, ErrDirIsFile
	}

	repo, err := git.PlainOpen(configFolder)
	if err != nil {
		if err == git.ErrRepositoryNotExists {
			return nil, err
		}

		return nil, err
	}

	return repo, nil
}
