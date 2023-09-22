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

type Repository struct {
	availableAuth []transport.AuthMethod
	auth          transport.AuthMethod
	repository    *git.Repository
	repoPath      string
	url           string
	author        *object.Signature
}

// NewRepository creates a new repository struct by either cloning or opening the repository
func NewRepository(config Git, repoPath string, clone bool) (*Repository, error) {
	var auth []transport.AuthMethod = nil
	if config.Repository == "" {
		return nil, errors.New("a repository URI is needed (either using GIT protocol or HTTPS)")
	}

	repoPath, err := ResolvePath(repoPath)
	if err != nil {
		return nil, err
	}

	for i, c := range config.SshAuth {
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
		availableAuth: auth,
		auth:          nil,
		repository:    nil,
		repoPath:      repoPath,
		url:           config.Repository,
		author: &object.Signature{
			Name:  config.Name,
			Email: config.Email,
		},
	}

	if clone {
		if err := repo.cloneRepository(); err != nil {
			return nil, err
		}
	} else {
		if err := repo.openRepository(); err != nil {
			return nil, err
		}
	}

	return repo, nil
}

// openRepository opens the repository at the given path
func (r *Repository) openRepository() error {
	repo, err := git.PlainOpen(r.repoPath)
	if err != nil {
		return err
	}

	w, err := repo.Worktree()
	if err != nil {
		return err
	}

	if r.availableAuth != nil {
		pulled := false
		for _, auth := range r.availableAuth {
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
			r.auth = auth
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

// PushChanges pushes changes to the remote repository
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

	author := r.author
	author.When = time.Now()
	if _, err := w.Commit(msg, &git.CommitOptions{
		Author: author,
	}); err != nil {
		return err
	}

	return r.repository.Push(&git.PushOptions{
		Auth: r.auth,
	})
}

// FetchChanges fetches changes from the remote repository. If changes are fetched,
// the function returns true, otherwise false.
func (r *Repository) FetchChanges() (bool, error) {
	if err := r.repository.Fetch(&git.FetchOptions{
		Auth: r.auth,
	}); err != nil {
		if err == git.NoErrAlreadyUpToDate {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

// cloneRepository clones the repository into the given path
func (r *Repository) cloneRepository() error {
	var repo *git.Repository
	var err error
	if r.availableAuth != nil {
		for _, auth := range r.availableAuth {
			repo, err = git.PlainClone(r.repoPath, false, &git.CloneOptions{
				URL:  r.url,
				Auth: auth,
			})
			if err != nil {
				if checkAuthErr(err) {
					logrus.WithField("auth", auth.String()).Warn("failed to authenticate. Trying next auth if exists")
					continue
				}

				return err
			}

			r.auth = auth
			break
		}

		if repo == nil {
			return fmt.Errorf("authentication failed for git repository")
		}
	} else {
		repo, err = git.PlainClone(r.repoPath, false, &git.CloneOptions{
			URL: r.url,
		})
		if err != nil {
			return err
		}
	}

	r.repository = repo

	return nil
}

// Pull pulls changes from the remote repository
func (r *Repository) Pull() error {
	w, err := r.repository.Worktree()
	if err != nil {
		return err
	}

	return w.Pull(&git.PullOptions{
		Auth: r.auth,
	})
}

// getSSHAuthMethod returns an authentication method for SSH
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

// checkAuthErr checks if the error is an authentication error
func checkAuthErr(err error) bool {
	return err == transport.ErrAuthorizationFailed || err == transport.ErrAuthenticationRequired
}
