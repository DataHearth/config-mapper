package git

import (
	"errors"
	"os"
	"time"

	"github.com/datahearth/config-mapper/internal/configuration"
	"github.com/datahearth/config-mapper/internal/misc"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
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

	if config.SSH.Passphrase != "" && config.SSH.PrivateKey != "" {
		privateKey, err := misc.AbsolutePath(config.SSH.PrivateKey)
		if err != nil {
			return nil, err
		}

		if _, err := os.Stat(privateKey); err != nil {
			return nil, err
		}

		auth, err = ssh.NewPublicKeysFromFile("git", privateKey, config.SSH.Passphrase)
		if err != nil {
			return nil, err
		}
	} else {
		auth = &http.BasicAuth{
			Username: config.BasicAuth.Username,
			Password: config.BasicAuth.Password,
		}
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
