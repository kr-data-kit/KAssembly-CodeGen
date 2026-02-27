package main

import (
	"flag"
	"fmt"
	"openassemblybinder/internal/command"
	"os"
)

func main() {
	fmt.Println("Open Assembly Client Code Generator v0.1.0")
	if len(os.Args) < 2 {
		fmt.Println("Invalid arguments.") // TODO: add help message
		return
	}

	subCommand := os.Args[1]

	switch subCommand {
	case "list":
		listCmd := flag.NewFlagSet("list", flag.ContinueOnError)
		method := listCmd.String("method", "simple", "Method to list services: simple or detailed")
		listCmd.Parse(os.Args[2:])
		fmt.Println("This feature is still under development")
		err := command.ListCommand(*method)
		if err != nil {
			panic(err)
		}
		return
	case "generate":
		genCmd := flag.NewFlagSet("generate", flag.ExitOnError)
		PackageName := genCmd.String("package", "openassembly", "Package name for generated code")
		ClientName := genCmd.String("client", "Client", "Client struct name for generated code")
		OutputPath := genCmd.String("output", "./generated", "Output path for generated code")
		CreateDir := genCmd.Bool("create-dir", false, "Create output directory if it does not exist")

		genCmd.Parse(os.Args[2:])
		err := command.GenerateCommand(
			*PackageName,
			*ClientName,
			*OutputPath,
			*CreateDir,
		)
		if err != nil {
			panic(err)
		}
		return
	default:
		fmt.Println("Unknown command:", subCommand)
		return
	}
}
