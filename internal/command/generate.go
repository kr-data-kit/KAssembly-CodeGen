package command

import (
	"fmt"
	"log/slog"
	"openassemblybinder/internal/generator"
	"os"

	"github.com/spf13/cobra"
)

var (
	language    string
	packageName string
	clientName  string
	outputPath  string
	createDir   bool
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate OpenAssembly client code",
	Long: `Generate OpenAssembly client code for Go and Python.

This command fetches the OpenAssembly API specification and generates
client code for the specified language and package.

Example:
  openassemblybinder generate --package myauth --output ./out
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if language == "" {
			return fmt.Errorf("language is required, use -m or --language (go or python)")
		}
		if language != "go" && language != "python" {
			return fmt.Errorf("unsupported language: %s (supported: go, python)", language)
		}

		slog.Info("Generating code with the following parameters")
		slog.Info("Generate parameters", "language", language, "package", packageName, "client", clientName, "output", outputPath, "create_dir", createDir)

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

		// Generate for selected language
		switch language {
		case "go":
			err := generator.GenerateGo(packageName, clientName, outputPath, createDir)
			if err != nil {
				return fmt.Errorf("Go code generation failed: %v", err)
			}
		case "python":
			err := generator.GeneratePython(packageName, outputPath, createDir)
			if err != nil {
				return fmt.Errorf("Python code generation failed: %v", err)
			}
		}

		slog.Info("Code generation completed successfully")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)

	generateCmdFlags := generateCmd.Flags()
	generateCmdFlags.StringVarP(&language, "language", "m", "", "Programming language (go, python) - REQUIRED")
	generateCmdFlags.StringVar(&packageName, "package", "openassemblyclient", "Package name for generated code")
	generateCmdFlags.StringVar(&clientName, "client", "OpenAssemblyClient", "Client struct name for generated code")
	generateCmdFlags.StringVar(&outputPath, "output", "./out", "Output path for generated code")
	generateCmdFlags.BoolVar(&createDir, "create-dir", true, "Create output directory if it does not exist")
}
