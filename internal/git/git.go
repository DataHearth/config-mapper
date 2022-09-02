package git

import (
	"errors"
	"fmt"
	"os"
	"time"

	"gitea.antoine-langlois.net/datahearth/config-mapper/internal/configuration"
	"gitea.antoine-langlois.net/datahearth/config-mapper/internal/misc"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/mitchellh/mapstructure"
)

var (
	ErrDirIsFile = errors.New("path is a file")
)

type RepositoryActions interface {
	PushChanges(msg string, newLines, removedLines []string) error
	GetWorktree() (*git.Worktree, error)
	GetAuthor() *object.Signature
	openRepository() error
}

type Repository struct {
	auth       transport.AuthMethod
	repository *git.Repository
	repoPath   string
	author     author
	url        string
}

type author struct {
	name  string
	email string
}

func NewRepository(config configuration.Git, repoPath string) (RepositoryActions, error) {
	var auth transport.AuthMethod
	if config.URL == "" {
		return nil, errors.New("a repository URI is needed (either using GIT protocol or HTTPS)")
	}
	repoPath, err := misc.AbsolutePath(repoPath)
	if err != nil {
		return nil, err
	}

	switch sshConfig := config.SSH.(type) {
	case map[string]interface{}:
		var outConfig configuration.Ssh
		if err := mapstructure.Decode(sshConfig, &outConfig); err != nil {
			return nil, err
		}

		auth, err = getSSHAuthMethod(outConfig)
		if err != nil {
			return nil, err
		}
	case []interface{}:
		for i, c := range sshConfig {
			if _, ok := c.(map[interface{}]interface{}); !ok {
				fmt.Printf("invalid format for configuration n°%d", i)
				continue
			}

			var outConfig configuration.Ssh
			if err := mapstructure.Decode(c, &outConfig); err != nil {
				fmt.Printf("failed to decode ssh configuration n°%d: %v\n", i, err)
				continue
			}

			auth, err = getSSHAuthMethod(outConfig)
			if err != nil {
				fmt.Printf("failed to create SSH authentication method for configuration n°%d: %v\n", i, err)
				continue
			}
		}

		if auth == nil {
			auth = &http.BasicAuth{
				Username: config.BasicAuth.Username,
				Password: config.BasicAuth.Password,
			}
		}
	default:
		return nil, errors.New("git ssh configuration canno't be unmarshaled. Please, pass a valid configuration")
	}

	repo := &Repository{
		auth:       auth,
		repository: nil,
		repoPath:   repoPath,
		url:        config.URL,
		author: author{
			name:  config.Name,
			email: config.Email,
		},
	}

	if err := repo.openRepository(); err != nil {
		return nil, err
	}

	return repo, nil
}

func (r *Repository) openRepository() error {
	s, err := os.Stat(r.repoPath)
	if err != nil {
		if os.IsNotExist(err) {
			repo, err := git.PlainClone(r.repoPath, false, &git.CloneOptions{
				URL:      r.url,
				Progress: os.Stdout,
				Auth:     r.auth,
			})
			if err != nil {
				return err
			}

			r.repository = repo
			return nil
		}

		return err
	}

	if !s.IsDir() {
		return ErrDirIsFile
	}

	repo, err := git.PlainOpen(r.repoPath)
	if err != nil {
		return err
	}

	w, err := repo.Worktree()
	if err != nil {
		return err
	}

	err = w.Pull(&git.PullOptions{
		Auth: r.auth,
	})
	if err != nil && err != git.NoErrAlreadyUpToDate {
		return err
	}

	r.repository = repo
	return nil
}

func (r *Repository) PushChanges(msg string, newLines, removedLines []string) error {
	w, err := r.repository.Worktree()
	if err != nil {
		return err
	}

	status, err := w.Status()
	if err != nil {
		return err
	}

	for file := range status {
		_, err = w.Add(file)
		if err != nil {
			return err
		}
	}

	if _, err := w.Commit(msg, &git.CommitOptions{
		Author: r.GetAuthor(),
	}); err != nil {
		return err
	}

	return r.repository.Push(&git.PushOptions{})
}

func (r *Repository) GetWorktree() (*git.Worktree, error) {
	return r.repository.Worktree()
}

func (r *Repository) GetAuthor() *object.Signature {
	return &object.Signature{
		Name:  r.author.name,
		Email: r.author.email,
		When:  time.Now(),
	}
}

func getSSHAuthMethod(config configuration.Ssh) (transport.AuthMethod, error) {
	if config.Passphrase == "" && config.PrivateKey == "" {
		return nil, errors.New("passphrase and private are empty")
	}

	privateKey, err := misc.AbsolutePath(config.PrivateKey)
	if err != nil {
		return nil, err
	}

	if _, err := os.Stat(privateKey); err != nil {
		return nil, err
	}

	auth, err := ssh.NewPublicKeysFromFile("git", privateKey, config.Passphrase)
	if err != nil {
		return nil, err
	}

	return auth, nil
}
