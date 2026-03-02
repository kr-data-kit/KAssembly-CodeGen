package command

import (
	"fmt"
	"openassemblybinder/internal/generator"
	"os"
)

func GenerateCommand(
	packageName string,
	clientName string,
	outputPath string,
	createDir bool,
) error {
	// just for test, print the parameters
	fmt.Println("Generating code with the following parameters:")
	fmt.Println("Package Name:", packageName)
	fmt.Println("Client Name :", clientName)
	fmt.Println("Output Path :", outputPath)
	fmt.Println("Create Dir  :", createDir)
	fmt.Println("")

	// TODO : add checking [y/n] before proceeding with code generation

	// check output directory existence (not creating directory here)
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		if createDir {
			err := os.MkdirAll(outputPath, os.ModePerm)
			if err != nil {
				return fmt.Errorf("failed to create output directory: %v", err)
			}
		} else {
			return fmt.Errorf("output directory does not exist: %s", outputPath)
		}
	}

	language := "Python" // for test

	switch language {
	case "Go":
		err := generator.GenerateGo(packageName, clientName, outputPath, createDir)
		if err != nil {
			return fmt.Errorf("code generation failed: %v", err)
		}
	case "Python":
		err := generator.GeneratePython(packageName, outputPath, createDir)
		if err != nil {
			return fmt.Errorf("code generation failed: %v", err)
		}
	default:
		return fmt.Errorf("unsupported language: %s", language)
	}

	fmt.Println("Code generation completed successfully.")
	return nil
}
