package mapper

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/datahearth/config-mapper/internal/configuration"
	"github.com/fatih/color"
	"github.com/gernest/wow"
	"github.com/gernest/wow/spin"
	"github.com/spf13/viper"
)

// InstallPackages install all packages from the configuration file by installation order
func InstallPackages(c configuration.PkgManagers) error {
	color.Blue("\n# Installing packages")

	for _, pkgManager := range c.InstallationOrder {
		var pkgs []string
		switch pkgManager {
		case "brew":
			pkgs = c.Brew
		case "apt":
			pkgs = c.Apt
		case "cargo":
			pkgs = c.Cargo
		case "npm":
			pkgs = c.Npm
		case "pip":
			pkgs = c.Pip
		case "go":
			pkgs = c.Go
		default:
			PrintError("package manager not supported: %s", pkgManager)
			continue
		}

		if _, err := exec.LookPath(pkgManager); err != nil {
			if pkgManager == "pip" {
				if _, err := exec.LookPath("pip3"); err != nil {
					return fmt.Errorf("%s and pip3 are not available on your system", pkgManager)
				}
				pkgManager = "pip3"
			} else {
				return fmt.Errorf("%s is not available on your system", pkgManager)
			}
		}

		if len(pkgs) == 0 {
			fmt.Printf("%s: nothing to do\n", pkgManager)
			return nil
		}

		cmd := exec.Command(pkgManager, "install")
		cmd.Args = append(cmd.Args, pkgs...)
		color.Blue("\n## Installing %s packages", pkgManager)

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
			PrintError("\n%s command failed: %v", pkgManager, err)
			return err
		}

		if v {
			// todo: find a way to clear spinner when done
			spinner.Stop()
		}
		color.Green("\n%s Packages intalled succesfully !", pkgManager)
	}

	return nil
}
