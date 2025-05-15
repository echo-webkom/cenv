package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/echo-webkom/cenv/internal"
	"github.com/fatih/color"
)

func showHelp() {
	fmt.Println("cenv [command] <args>")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("    check     Check if .env matches schema")
	fmt.Println("    update    Generate schema based on .env")
	fmt.Println("    fix       Fix will")
	fmt.Println("                 create .env if one does not exist")
	fmt.Println("                 fill in default/public values for empty fields")
	fmt.Println("                 automatically fix issues with the .env if any")
	fmt.Println("                 always reuse values already in the .env")
	fmt.Println("    tags      Show list of available tags and their uses")
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

func showTags() {
	fmt.Println("List of tags")
	fmt.Println()
	fmt.Println("   required        Field has to have a non-empty value")
	fmt.Println("   public          Field is static and shown in schema")
	fmt.Println("   default <value> Field will be set to given default value if empty")
	fmt.Println("   length <n>      Field must have given length")
	fmt.Println("   format <fmt>    Field must have given format (gokenizer pattern)")
	fmt.Println("                   See github.com/jesperkha/gokenizer")
	fmt.Println("   enum <values>   Field must be one of the given values")
	fmt.Println("                   Values are separated by a |")
	fmt.Println("                   Example: @enum user | admin | guest")
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

	if command == "tags" {
		showTags()
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
		if err := internal.Fix(envPath, schemaPath); err != nil {
			errorExit(err)
		}
		return
	}

	if command == "check" {
		if err := internal.Check(envPath, schemaPath); err != nil {
			errorExit(err)
		}
		return
	}

	if command == "update" {
		if err := internal.Update(envPath, schemaPath); err != nil {
			errorExit(err)
		}
		return
	}

	errorExitS(fmt.Sprintf("'%s' is not a command", command))
}
