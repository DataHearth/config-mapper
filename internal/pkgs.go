package mapper

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

var (
	errLogger             = log.New(os.Stderr, "", 0)
	ErrFailedInstallation = errors.New("failed to install some packages. Please, checkout STDERR for more information")
	ErrFailedSaving       = errors.New("failed to save some packages. Please, checkout STDERR for more information")
	ErrBrewNotAvailable   = errors.New("homebrew is not available on your system")
	ErrAptNotAvailable    = errors.New("aptitude is not available on your system")
)

func LoadPkgs(c PkgManagers) error {
	color.Blue("# Load folders into saved location")

	for _, pkg := range c.InstallationOrder {
		switch pkg {
		case "homebrew":
			if err := installBrewPkgs(c.Homebrew); err != nil {
				PrintError(err.Error())
				return ErrFailedInstallation
			}
		case "apt":
			if err := installAptPkgs(c.Aptitude); err != nil {
				PrintError(err.Error())
				return ErrFailedInstallation
			}
		}
	}

	return nil
}

func SavePkgs(cfg Configuration) error {
	color.Blue("# Save user installed packages")

	for _, pkg := range cfg.PackageManagers.InstallationOrder {
		switch pkg {
		case "homebrew":
			if err := SaveBrewPkgs(cfg); err != nil {
				PrintError(err.Error())
				return ErrFailedSaving
			}
		case "apt":
			fmt.Println("implemented soon!")
		}
	}

	return nil
}

func SaveBrewPkgs(cfg Configuration) error {
	if _, err := exec.LookPath("brew"); err != nil {
		return err
	}

	color.Blue("## Saving Homebrew packages")

	o, err := exec.Command("brew", "leaves", "--installed-on-request").Output()
	if err != nil {
		return err
	}

	pkgs := strings.Split(string(o), "\n")
	cfg.PackageManagers.Homebrew = append([]string{}, pkgs[:len(pkgs)-1]...)

	b, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	if err := os.WriteFile(viper.GetString("configuration-file"), b, 0755); err != nil {
		return err
	}

	color.Green("Packages saved succesfully !")
	return nil
}

func installBrewPkgs(pkgs []string) error {
	if _, err := exec.LookPath("brew"); err != nil {
		return ErrBrewNotAvailable
	}

	if len(pkgs) == 0 {
		fmt.Println("homebrew: nothing to do")
		return nil
	}

	cmd := exec.Command("brew", "install")
	cmd.Args = append(cmd.Args, pkgs...)
	color.Blue("## Installing Homebrew packages")

	if err := cmd.Run(); err != nil {
		PrintError("brew command failed: %v", err)
		return err
	}

	color.Green("Packages intalled succesfully !")

	return nil
}

func installAptPkgs(pkgs []string) error {
	if _, err := exec.LookPath("apt-get"); err != nil {
		return ErrAptNotAvailable
	}

	if len(pkgs) == 0 {
		fmt.Println("aptitude: nothing to do")
		return nil
	}

	cmd := exec.Command("sudo", "apt-get", "install")
	cmd.Args = append(cmd.Args, pkgs...)

	color.Blue("## Installing aptitude packages")

	if err := cmd.Run(); err != nil {
		PrintError("aptitude command failed: %v", err)
		return err
	}

	color.Green("Packages intalled succesfully !")

	return nil
}
