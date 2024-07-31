package mapper

import (
	"errors"
	"log"
	"os"
	"os/exec"

	"github.com/pterm/pterm"
)

var (
	errLogger             = log.New(os.Stderr, "", 0)
	ErrFailedInstallation = errors.New("failed to installed some packages. Please, checkout STDERR for more information")
)

func LoadPkgs(c PkgManagers) error {
	for _, pkg := range c.InstallationOrder {
		switch pkg {
		case "homebrew":
			if err := installBrewPkgs(c.Homebrew); err != nil {
				return ErrFailedInstallation
			}
		case "apt":
			if err := installAptPkgs(c.Aptitude); err != nil {
				return ErrFailedInstallation
			}
		}
	}

	return nil
}

func installBrewPkgs(pkgs []string) error {
	if _, err := exec.LookPath("brew"); err != nil {
		errLogger.Println(pterm.Red("Homebrew is not installed on your system"))
		return nil
	}

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

	introSpinner.Success("Packages intalled succesfully")

	return nil
}

func installAptPkgs(pkgs []string) error {
	if _, err := exec.LookPath("apt-get"); err != nil {
		errLogger.Println(pterm.Red("aptitude is not available on your system"))
		return nil
	}

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

	introSpinner.Success("Packages intalled succesfully")

	return nil
}
