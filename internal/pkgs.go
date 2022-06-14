package mapper

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/datahearth/config-mapper/internal/configuration"
	"github.com/fatih/color"
	"github.com/gernest/wow"
	"github.com/gernest/wow/spin"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

var (
	ErrFailedInstallation = errors.New("failed to install some packages. Please, checkout STDERR for more information")
	ErrFailedSaving       = errors.New("failed to save some packages. Please, checkout STDERR for more information")
	ErrBrewNotAvailable   = errors.New("homebrew is not available on your system")
	ErrAptNotAvailable    = errors.New("aptitude is not available on your system")
)

// LoadPkgs triggers related functions with passed order
func LoadPkgs(c configuration.PkgManagers) error {
	color.Blue("\n# Installing packages")

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

// SavePkgs triggers related functions with passed order
func SavePkgs(cfg configuration.Configuration) error {
	color.Blue("# Saving user installed packages")

	for _, pkg := range cfg.PackageManagers.InstallationOrder {
		switch pkg {
		case "homebrew":
			if err := saveBrewPkgs(cfg); err != nil {
				PrintError(err.Error())
				return ErrFailedSaving
			}
		case "apt":
			fmt.Println("implemented soon!")
		}
	}

	return nil
}

// saveBrewPkgs gather user installed packages by running `brew leaves --installed-on-request`.
// It captures the output, parse it and save it into the configuration.
func saveBrewPkgs(cfg configuration.Configuration) error {
	if _, err := exec.LookPath("brew"); err != nil {
		return err
	}

	color.Blue("\n## Saving Homebrew packages")

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

// installBrewPkgs installs homebrew packages by passing them to homebrew's CLI.
// STDERR and STDOUT are captured if verbose flag is passed.
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
	color.Blue("\n## Installing Homebrew packages")

	spinner := wow.New(os.Stdout, spin.Get(spin.Dots3), " Running...")

	v := viper.GetBool("verbose")
	if v {
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
	} else {
		spinner.Start()
	}

	if err := cmd.Run(); err != nil {
		spinner.Stop()
		PrintError("brew command failed: %v", err)
		return err
	}

	if v {
		// todo: find a way to clear spinner when done
		spinner.Stop()
	}
	color.Green("\nPackages intalled succesfully !")

	return nil
}

// installAptPkgs installs all provided "apt" packages by passing them to the Advanced Package Tool's CLI
func installAptPkgs(pkgs []string) error {
	if _, err := exec.LookPath("apt"); err != nil {
		return ErrAptNotAvailable
	}

	if len(pkgs) == 0 {
		fmt.Println("apt: nothing to do")
		return nil
	}

	cmd := exec.Command("sudo", "apt", "install")
	cmd.Args = append(cmd.Args, pkgs...)

	color.Blue("\n## Installing apt packages")

	if err := cmd.Run(); err != nil {
		PrintError("apt command failed: %v", err)
		return err
	}

	color.Green("Packages intalled succesfully !")

	return nil
}
