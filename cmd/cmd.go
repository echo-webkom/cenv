package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/echo-webkom/cenv/cenv"
	"github.com/fatih/color"
)

func showHelp() {
	fmt.Println("cenv [command] <args>")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("    check     Check if .env matches schema")
	fmt.Println("    update    Generate schema based on .env")
	fmt.Println("    fix       Automatically fix issues with the .env")
	fmt.Println("              Tries to reuse previous env values")
	fmt.Println("    help      Show this help message")
	fmt.Println("    version   Show version")
	fmt.Println("    upgrade   Upgrade to latest version")
	fmt.Println()
	fmt.Println("Flags:")
	fmt.Println("    --env <path>      Path to env file, default is current dir")
	fmt.Println("    --schema <path>   Path to schema file, default is current dir")
	fmt.Println("    --skip-version    Skip version check")
	fmt.Println("    --version         Show version")
	fmt.Println()
}

func errorExitS(message string) {
	color.RGB(237, 93, 83).Println(fmt.Sprintf("cenv: %s", message))
	os.Exit(1)
}

func errorExit(err error) {
	color.RGB(237, 93, 83).Println(err.Error())
	os.Exit(1)
}

func Run() {
	checkIfLatestVersion()

	if len(os.Args) < 2 {
		showHelp()
		os.Exit(1)
	}

	envPath := ".env"
	schemaPath := "cenv.schema.json"
	command := os.Args[1]

	i := 2
	for i < len(os.Args) {
		arg := os.Args[i]

		if len(os.Args) <= i+1 {
			errorExitS("expected value after flag " + arg)
		}

		val := os.Args[i+1]

		switch arg {
		case "--env":
			envPath = val
		case "--schema":
			schemaPath = val
		default:
			errorExitS("unknown flag " + arg)
		}

		i += 2
	}

	if command == "help" || command == "-h" || command == "--help" {
		showHelp()
		return
	}

	if command == "version" || command == "--version" || command == "-v" {
		if Version == "" {
			fmt.Println("you are running a development version")
		} else {
			fmt.Println(Version)
		}
		return
	}

	if command == "upgrade" {
		cmd := exec.Command("bash", "-c", "curl -fsSL https://raw.githubusercontent.com/echo-webkom/cenv/refs/heads/main/install.sh | bash")
		fmt.Println("Installing latest release...")

		if _, err := cmd.Output(); err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("Upgrade complete")
		return
	}

	if command == "fix" {
		if err := cenv.Fix(envPath, schemaPath); err != nil {
			errorExit(err)
		}
		return
	}

	if command == "check" {
		if err := cenv.Check(envPath, schemaPath); err != nil {
			errorExit(err)
		}
		return
	}

	if command == "update" {
		if err := cenv.Update(envPath, schemaPath); err != nil {
			errorExit(err)
		}
		return
	}

	errorExitS(fmt.Sprintf("'%s' is not a command", command))
}
