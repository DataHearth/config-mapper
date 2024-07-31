package mapper

import (
	"errors"
	"log"
	"os"
	"os/exec"

	"github.com/pterm/pterm"
	"github.com/spf13/viper"
)

var (
	errLogger             = log.New(os.Stderr, "", 0)
	ErrFailedInstallation = errors.New("failed to installed some packages. Please, checkout STDERR for more information")
)

func LoadPkgs() error {
	order := viper.GetStringSlice("package-managers.installation-order")

	for _, pkg := range order {
		switch pkg {
		case "homebrew":
			if err := installBrewPkgs(); err != nil {
				return ErrFailedInstallation
			}
		case "apt":
			if err := installAptPkgs(); err != nil {
				return ErrFailedInstallation
			}
		}
	}

	return nil
}

func installBrewPkgs() error {
	if _, err := exec.LookPath("brew"); err != nil {
		errLogger.Println(pterm.Red("Homebrew is not installed on your system"))
	}

	pkgs := viper.GetStringSlice("package-managers.homebrew")
	if len(pkgs) == 0 {
		pterm.Println(pterm.Blue("homebrew: nothing to do"))
		return nil
	}

	cmd := exec.Command("brew", "install")
	cmd.Args = append(cmd.Args, pkgs...)
	introSpinner, _ := pterm.DefaultSpinner.WithShowTimer(true).WithRemoveWhenDone(false).Start("Installing homebrew packages")

	chErr := make(chan error)
	defer close(chErr)

	go func(chErr chan error) {
		if err := cmd.Run(); err != nil {
			chErr <- err
			return
		}

		chErr <- nil
	}(chErr)

	err := <-chErr
	introSpinner.Stop()
	if err != nil {
		return err
	}

	return nil
}

func installAptPkgs() error {
	if _, err := exec.LookPath("apt-get"); err != nil {
		errLogger.Println(pterm.Red("aptitude is not available on your system"))
	}

	pkgs := viper.GetStringSlice("package-managers.apt")
	if len(pkgs) == 0 {
		pterm.Println(pterm.Blue("aptitude: nothing to do"))
		return nil
	}

	cmd := exec.Command("sudo", "apt-get", "install")
	cmd.Args = append(cmd.Args, pkgs...)

	introSpinner, _ := pterm.DefaultSpinner.WithShowTimer(true).WithRemoveWhenDone(false).Start("Installing aptitude packages")

	chErr := make(chan error)
	defer close(chErr)

	go func(chErr chan error) {
		if err := cmd.Run(); err != nil {
			chErr <- err
			return
		}

		chErr <- nil
	}(chErr)

	err := <-chErr
	introSpinner.Stop()
	if err != nil {
		return err
	}

	return nil
}
