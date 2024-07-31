package mapper

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"gitea.antoine-langlois.net/datahearth/config-mapper/internal/configuration"
	"github.com/charmbracelet/log"
	"github.com/gernest/wow"
	"github.com/gernest/wow/spin"
	"github.com/spf13/viper"
)

// InstallPackages install all packages from the configuration file by installation order
func InstallPackages(c configuration.PkgManagers) error {
	pkgManagers := map[string]bool{}
	for _, pkgManager := range viper.GetStringSlice("exclude-pkg-managers") {
		pkgManagers[pkgManager] = true
	}

	for _, pkgManager := range c.InstallationOrder {
		log.Info("installing packages", "package-manager", pkgManager)
		if _, ok := pkgManagers[pkgManager]; ok {
			log.Info("skipping package manager", "package-manager", pkgManager)
			continue
		}

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
		case "nala":
			pkgs = c.Nala
		default:
			log.Error("package manager not supported", "package-manager", pkgManager)
			continue
		}

		if _, err := exec.LookPath(pkgManager); err != nil {
			// * pip might not be available on the system but pip3 is
			if pkgManager == "pip" {
				if _, err := exec.LookPath("pip3"); err != nil {
					log.Error("pip and pip3 are not available on your system", "package-manager", pkgManager)
					continue
				}
				pkgManager = "pip3"
			} else {
				log.Error("package manager not available on your system", "package-manager", pkgManager)
				continue
			}
		}
		// * for some reason, apt binary is available on darwin. Exclude it to avoid errors
		if pkgManager == "apt" && runtime.GOOS == "darwin" {
			log.Error("package manager not available on your system", "package-manager", pkgManager)
			continue
		}

		if len(pkgs) == 0 {
			fmt.Printf("✔️ nothing to do\n\n")
			continue
		}

		v := viper.GetBool("verbose")
		commands := []*exec.Cmd{}
		// * package managers requiring sudo permission
		if pkgManager == "apt" || pkgManager == "nala" {
			commands = append(commands, buildDefaultCommand([]string{"sudo", pkgManager, "install", "-y"}, pkgs, v))
		} else if pkgManager == "cargo" {
			commands = buildCargoCommand(pkgs, v)
		} else {
			commands = append(commands, buildDefaultCommand([]string{pkgManager, "install"}, pkgs, v))
		}

		for i, cmd := range commands {
			spinner := wow.New(os.Stdout, spin.Get(spin.Dots3), " Installing...")
			if !v {
				spinner.Start()
			}
			if err := cmd.Run(); err != nil {
				if v {
					log.Error(err)
				} else {
					msg := fmt.Sprintf(" %s", cmd.Args)
					if i == len(commands)-1 {
						msg = fmt.Sprintf("%s\n", msg)
					}
					spinner.PersistWith(spin.Spinner{Frames: []string{"❌"}}, msg)
				}
				continue
			}

			if !v {
				// msg := fmt.Sprintf(" %s %s", color.GreenString("Success\t"), cmd.Args)
				msg := fmt.Sprintf(" %s", cmd.Args)
				if i == len(commands)-1 {
					msg = fmt.Sprintf("%s\n", msg)
				}
				spinner.PersistWith(spin.Spinner{Frames: []string{"✔️"}}, msg)
			}
		}
	}

	return nil
}

func buildCargoCommand(packages []string, verbose bool) []*exec.Cmd {
	commands := []*exec.Cmd{}

	cmd := exec.Command("cargo", "install")
	for _, pkg := range packages {
		if strings.Contains(pkg, " ") {
			customCmd := exec.Command("cargo", "install")
			customCmd.Args = append(cmd.Args, strings.Split(pkg, " ")...)
			if verbose {
				customCmd.Stderr = os.Stderr
				customCmd.Stdout = os.Stdout
			}
			commands = append(commands, customCmd)
		} else {
			cmd.Args = append(cmd.Args, pkg)
		}
	}

	if verbose {
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
	}
	if len(cmd.Args) > 2 {
		commands = append(commands, cmd)
	}

	return commands
}

func buildDefaultCommand(command, packages []string, verbose bool) *exec.Cmd {
	cmd := exec.Command(command[0], command[1:]...)
	cmd.Args = append(cmd.Args, packages...)
	if verbose {
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
	}

	return cmd
}
