package configuration

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	osUser "os/user"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
)

func loadConfigSSH(uri string) error {
	config, host, path, err := getSSHConfig(uri)
	if err != nil {
		return err
	}

	c, err := ssh.Dial("tcp", host, config)
	if err != nil {
		return err
	}
	defer c.Close()

	s, err := c.NewSession()
	if err != nil {
		return err
	}
	defer s.Close()

	buff := new(bytes.Buffer)
	s.Stdout = buff
	if err := s.Run(fmt.Sprintf("cat %s", path)); err != nil {
		return err
	}

	if err := viper.ReadConfig(buff); err != nil {
		return err
	}

	return nil
}

func getSSHConfig(uriFlag string) (*ssh.ClientConfig, string, string, error) {
	var err error
	var user, passwd, host, configPath, key string
	uri := strings.Split(uriFlag, "ssh://")[1]

	if key = viper.GetString("ssh-key"); key != "" {
		uri, user, passwd, err = getCredentials(uri)
		if err != nil {
			return nil, "", "", err
		}

		host, configPath, err = getUriContent(uri)
		if err != nil {
			return nil, "", "", err
		}
	} else if user = viper.GetString("ssh-user"); user != "" {
		host, configPath, err = getUriContent(uri)
		if err != nil {
			return nil, "", "", err
		}

		passwd = viper.GetString("ssh-password")
	} else {
		uri, user, passwd, err = getCredentials(uri)
		if err != nil {
			return nil, "", "", err
		}
		if passwd == "" {
			passwd = viper.GetString("ssh-password")
		}

		host, configPath, err = getUriContent(uri)
		if err != nil {
			return nil, "", "", err
		}
	}

	if user == "" {
		color.Yellow("WARNING: no user was found in either the URI and flags. Current user will be used")

		var currentUser *osUser.User
		currentUser, err = osUser.Current()
		if err != nil {
			return nil, "", "", err
		}

		user = currentUser.Username
	}

	var auth ssh.AuthMethod
	if key != "" {
		auth, err = createPubKeyAuth(key)
		if err != nil {
			return nil, "", "", err
		}
	} else {
		auth = ssh.Password(passwd)
	}

	h, err := os.UserHomeDir()
	if err != nil {
		return nil, "", "", err
	}

	hostKeyCallback, err := knownhosts.New(fmt.Sprintf("%s/.ssh/known_hosts", h))
	if err != nil {
		return nil, "", "", err
	}

	if len(strings.SplitN(host, ":", 1)) == 1 {
		host = fmt.Sprintf("%s:22", host)
	}

	return &ssh.ClientConfig{
		User:            user,
		Auth:            []ssh.AuthMethod{auth},
		HostKeyCallback: hostKeyCallback,
	}, host, configPath, nil
}

// getCredentials takes an SSH URI and returns (splitted URI, user, passwd, error)
//
// "passwd" can be empty in case of a single credential in URI
func getCredentials(uri string) (string, string, string, error) {
	uriContent := strings.SplitN(uri, "@", 2)
	if len(uriContent) == 1 {
		fmt.Printf("uriContent: %v\n", uriContent)
		return "", "", "", errors.New("no credentials in URI")
	}

	credentials := strings.SplitN(uriContent[0], ":", 2)
	if len(credentials) == 1 {
		return uriContent[1], credentials[0], "", nil
	}

	return uriContent[1], credentials[0], credentials[1], nil
}

// getUriContent takes an SSH URI and returns (host, path, error)
func getUriContent(uri string) (string, string, error) {
	uriContent := strings.Split(uri, ":")
	if len(uriContent) < 2 {
		return "", "", errors.New("ssh URI is malformed. It's missing either a host or path. E.g: \"ssh://localhost:/my/config/file.yml\"")
	}

	return uriContent[0], uriContent[1], nil
}

func createPubKeyAuth(key string) (ssh.AuthMethod, error) {
	var signer ssh.Signer
	privateKey, err := ioutil.ReadFile(key)
	if err != nil {
		return nil, err
	}

	if passphrase := os.Getenv("CONFIG_MAPPER_SSH_PASSPHRASE"); passphrase != "" {
		signer, err = ssh.ParsePrivateKeyWithPassphrase(privateKey, []byte(passphrase))
	} else {
		signer, err = ssh.ParsePrivateKey(privateKey)
	}

	if err != nil {
		return nil, err
	}

	return ssh.PublicKeys(signer), nil
}
