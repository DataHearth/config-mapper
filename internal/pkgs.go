package mapper

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

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
		color.Blue("## Installing %s packages", pkgManager)

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
			PrintError("package manager not supported: %s\n", pkgManager)
			continue
		}

		if _, err := exec.LookPath(pkgManager); err != nil {
			// * pip might not be available on the system but pip3 is
			if pkgManager == "pip" {
				if _, err := exec.LookPath("pip3"); err != nil {
					PrintError("%s and pip3 are not available on your system\n", pkgManager)
					continue
				}
				pkgManager = "pip3"
			} else {
				PrintError("%s is not available on your system\n", pkgManager)
				continue
			}
		}
		// * for some reason, apt binary is available on darwin. Exclude it to avoid errors
		if pkgManager == "apt" && runtime.GOOS == "darwin" {
			PrintError("%s is not available on your system\n", pkgManager)
			continue
		}

		if len(pkgs) == 0 {
			fmt.Printf("%s: nothing to do\n", pkgManager)
			continue
		}

		cmd := exec.Command(pkgManager, "install")
		for _, pkg := range pkgs {
			if strings.Contains(pkg, " ") {
				cmd.Args = append(cmd.Args, strings.Split(pkg, " ")...)
			} else {
				cmd.Args = append(cmd.Args, pkg)
			}
		}

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
			fmt.Println()
		}
		color.Green("%s packages intalled succesfully !", pkgManager)
		fmt.Println()
	}

	return nil
}
