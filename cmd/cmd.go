package cmd

import (
	"fmt"
	"os"

	"github.com/echo-webkom/cenv/cenv"
)

func showHelp() {
	fmt.Println("cenv [command] <args>")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("	check     Check if .env matches schema")
	fmt.Println("	update    Generate schema based on .env")
	fmt.Println()
	fmt.Println("Flags:")
	fmt.Println("	--env <path>	   Path to env file, default is current dir")
	fmt.Println("	--schema <path>    Path to schema file, default is current dir")
	fmt.Println()
}

func errorHelp(message string) {
	fmt.Printf("cenv: %s\n\n", message)
	showHelp()
	os.Exit(1)
}

func errorExit(message string) {
	fmt.Printf("cenv: %s\n", message)
	os.Exit(1)
}

func Run() {
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
			errorExit("expected value after flag " + arg)
		}

		if arg == "--env" {
			envPath = arg
		} else if arg == "--schema" {
			schemaPath = arg
		} else {
			errorExit("unknown flag " + arg)
		}

		i += 2
	}

	if command == "help" {
		showHelp()
		return
	}

	if command == "check" {
		if err := cenv.Check(envPath, schemaPath); err != nil {
			errorExit(err.Error())
		}
		return
	}

	if command == "update" {
		if err := cenv.Update(envPath, schemaPath); err != nil {
			errorExit(err.Error())
		}
		return
	}

	errorExit(fmt.Sprintf("'%s' is not a command", command))
}
