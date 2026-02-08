package main

import (
	"flag"
	"fmt"
	"openassemblybind/internal/command"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("Test")
	if len(os.Args) < 2 {
		fmt.Println("Invalid arguments.") // TODO: add help message
		return
	}

	// always load .env to get ASSEMBLY_API_KEY
	// TODO: .env path config
	if !checkKey(".env") {
		panic("ASSEMBLY_API_KEY environment variable is not set")
	}
	key := os.Getenv("ASSEMBLY_API_KEY")

	subCommand := os.Args[1]

	switch subCommand {
	case "list":
		listCmd := flag.NewFlagSet("list", flag.ContinueOnError)
		method := listCmd.String("method", "simple", "Method to list services: simple or detailed")
		listCmd.Parse(os.Args[2:])
		err := command.ListCommand(key, *method)
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
			key,
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

func checkKey(
	godotenvPath string,
) bool {
	if os.Getenv("ASSEMBLY_API_KEY") != "" {
		return true
	}
	err := godotenv.Load(godotenvPath)
	if err != nil {
		return false
	}
	if os.Getenv("ASSEMBLY_API_KEY") != "" {
		return true
	}
	return false
}
