package internal

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/sirupsen/logrus"
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
	auth       []transport.AuthMethod
	repository *git.Repository
	repoPath   string
	author     author
	url        string
}

type author struct {
	name  string
	email string
}

func NewRepository(config Git, repoPath string) (RepositoryActions, error) {
	var auth []transport.AuthMethod = nil
	if config.Repository == "" {
		return nil, errors.New("a repository URI is needed (either using GIT protocol or HTTPS)")
	}

	repoPath, err := ResolvePath(repoPath)
	if err != nil {
		return nil, err
	}

	for i, c := range config.SSH {
		sshAuth, err := getSSHAuthMethod(c)
		if err != nil {
			fmt.Printf("failed to create SSH authentication method for configuration n°%d: %v\n", i, err)
			continue
		}

		if auth == nil {
			auth = []transport.AuthMethod{}
		}
		auth = append(auth, sshAuth)
	}

	if len(auth) == 0 {
		auth = append(auth, &http.BasicAuth{
			Username: config.BasicAuth.Username,
			Password: config.BasicAuth.Password,
		})
	}

	repo := &Repository{
		auth:       auth,
		repository: nil,
		repoPath:   repoPath,
		url:        config.Repository,
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
			var repo *git.Repository
			var err error
			if r.auth != nil {
				for _, auth := range r.auth {
					repo, err = git.PlainClone(r.repoPath, false, &git.CloneOptions{
						URL:      r.url,
						Progress: os.Stdout,
						Auth:     auth,
					})
					if err != nil {
						if checkAuthErr(err) {
							logrus.WithField("auth", auth.String()).Warn("failed to authenticate. Trying next auth if exists")
							continue
						}

						return err
					}

					break
				}

				if repo == nil {
					return fmt.Errorf("authentication failed for git repository")
				}
			} else {
				repo, err = git.PlainClone(r.repoPath, false, &git.CloneOptions{
					URL:      r.url,
					Progress: os.Stdout,
				})
				if err != nil {
					return err
				}
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

	if r.auth != nil {
		pulled := false
		for _, auth := range r.auth {
			err := w.Pull(&git.PullOptions{
				Auth: auth,
			})
			if err != nil {
				if err == git.NoErrAlreadyUpToDate {
					pulled = true
					break
				} else if checkAuthErr(err) {
					logrus.WithField("auth", auth.String()).Warn("failed to authenticate. Trying next auth if exists")
					continue
				}

				return err
			}

			pulled = true
			break
		}

		if !pulled {
			return fmt.Errorf("authentication failed for git repository")
		}
	} else {
		if err := w.Pull(&git.PullOptions{}); err != nil && err != git.NoErrAlreadyUpToDate {
			return err
		}
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
		if _, err := w.Add(file); err != nil {
			return err
		}
	}

	if _, err := w.Commit(msg, &git.CommitOptions{
		Author: r.GetAuthor(),
	}); err != nil {
		return err
	}

	if r.auth != nil {
		pushed := false
		for _, auth := range r.auth {
			err := r.repository.Push(&git.PushOptions{
				Auth: auth,
			})
			if err != nil {
				if checkAuthErr(err) {
					logrus.WithField("auth", auth.String()).Warn("failed to authenticate. Trying next auth if exists")
					continue
				}

				return err
			}

			pushed = true
			break
		}

		if !pushed {
			return fmt.Errorf("authentication failed for git repository")
		}
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

func getSSHAuthMethod(config SshAuth) (transport.AuthMethod, error) {
	if config.PrivateKey == "" {
		return nil, errors.New("\"private-key\" field is empty")
	}

	privateKey, err := ResolvePath(config.PrivateKey)
	if err != nil {
		return nil, err
	}

	s, err := os.Stat(privateKey)
	if err != nil {
		return nil, err
	}
	if s.IsDir() {
		return nil, errors.New("private key is a directory")
	}

	auth, err := ssh.NewPublicKeysFromFile("git", privateKey, config.Passphrase)
	if err != nil {
		return nil, err
	}

	return auth, nil
}

func checkAuthErr(err error) bool {
	return err == transport.ErrAuthorizationFailed || err == transport.ErrAuthenticationRequired
}
