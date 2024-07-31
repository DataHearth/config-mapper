package mapper

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
)

var (
	ErrDirIsFile  = errors.New("path is a file")
	ErrInvalidEnv = errors.New("found invalid environment variable in path")
)

func OpenGitRepo(c Git, l string) (*git.Repository, error) {
	s, err := os.Stat(l)
	if err != nil {
		if os.IsNotExist(err) {
			var auth transport.AuthMethod

			if c.SSH.Passphrase != "" && c.SSH.PrivateKey != "" {
				privateKey, err := absolutePath(c.SSH.PrivateKey)
				if err != nil {
					return nil, err
				}

				if _, err := os.Stat(privateKey); err != nil {
					return nil, err
				}

				auth, err = ssh.NewPublicKeysFromFile("git", privateKey, c.SSH.Passphrase)
				if err != nil {
					return nil, err
				}
			} else {
				auth = &http.BasicAuth{
					Username: c.BasicAuth.Username,
					Password: c.BasicAuth.Password,
				}
			}

			repo, err := git.PlainClone(l, false, &git.CloneOptions{
				URL:      c.Repository,
				Progress: os.Stdout,
				Auth:     auth,
			})
			if err != nil {
				return nil, err
			}

			return repo, nil
		}

		return nil, err
	}

	if !s.IsDir() {
		return nil, ErrDirIsFile
	}

	repo, err := git.PlainOpen(l)
	if err != nil {
		return nil, err
	}

	return repo, nil
}

func absolutePath(p string) (string, error) {
	finalPath := p
	if strings.Contains(finalPath, "~") {
		h, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}

		finalPath = strings.Replace(p, "~", h, 1)
	}

	splitted := strings.Split(finalPath, "/")
	finalPath = ""
	for _, s := range splitted {
		pathPart := s
		if strings.Contains(s, "$") {
			env := os.Getenv(s)
			if env == "" {
				return "", ErrInvalidEnv
			}
			pathPart = env
		}

		finalPath += fmt.Sprintf("/%s", pathPart)
	}

	return path.Clean(finalPath), nil
}
