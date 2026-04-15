package command

import (
	"context"
	"fmt"
	"kassemblycodegen/internal/generator"
	gogen "kassemblycodegen/internal/generator/go"
	pygen "kassemblycodegen/internal/generator/python"
	"kassemblycodegen/internal/setting"
	"log/slog"
	"os"

	"github.com/spf13/cobra"
)

var (
	language         string
	packageName      string
	outputPath       string
	createDir        bool
	confirmYes       bool
	useCache         bool
	includeEndpoints []string
	excludeEndpoints []string
	goMod            bool
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate OpenAssembly client code",
	Long: `Generate OpenAssembly client code for Go and Python.

This command fetches the OpenAssembly API specification and generates
client code for the specified language and package.

Example:
  kassemblycodegen generate --package myauth --output ./out
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if language == "" {
			return fmt.Errorf("language is required, use -m or --language (go or python)")
		}
		if language != "go" && language != "python" {
			return fmt.Errorf("unsupported language: %s (supported: go, python)", language)
		}

		if err := validatePackageName(language, packageName); err != nil {
			return err
		}

		fmt.Printf("%s\n", "┌──────────────────────────────────────────────────────────┐")
		fmt.Printf("│ %-56s │\n", "Code Generation Configuration")
		fmt.Printf("%s\n", "├──────────────────────────────────────────────────────────┤")
		fmt.Printf("│ %-17s %-38s │\n", "Language:", language)
		fmt.Printf("│ %-17s %-38s │\n", "Package Name:", packageName)
		fmt.Printf("│ %-17s %-38s │\n", "Output Path:", outputPath)
		fmt.Printf("│ %-17s %-38v │\n", "Create Directory:", createDir)
		fmt.Printf("│ %-17s %-38v │\n", "Use Cache:", useCache)
		fmt.Printf("│ %-17s %-38v │\n", "Go Mod:", goMod)
		fmt.Printf("%s\n", "└──────────────────────────────────────────────────────────┘")

		if !confirmYes {
			fmt.Print("\nProceed with code generation? [y/N]: ")
			var response string
			fmt.Scanln(&response)
			if response != "y" && response != "Y" {
				fmt.Println("Generation cancelled.")
				return nil
			}
		}
		fmt.Println("Starting generation...")

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
		ctx, cancel := context.WithCancel(cmd.Context())
		defer cancel()

		source := newEndpointSource(useCache, includeEndpoints, excludeEndpoints)

		switch language {
		case "go":
			set := setting.GoSetting{
				GlobalSetting: setting.GlobalSetting{
					Path:      outputPath,
					CreateDir: createDir,
				},
				PackageName: packageName,
				IsMod:       goMod,
			}
			gen := gogen.NewClientGenerator(set)
			err := generator.Generate(ctx, gen, source)
			if err != nil {
				return fmt.Errorf("Go code generation failed: %v", err)
			}
		case "python":
			set := setting.PythonSetting{
				GlobalSetting: setting.GlobalSetting{
					Path:      outputPath,
					CreateDir: createDir,
				},
				PackageName: packageName,
			}
			gen := pygen.NewClientGenerator(set)
			err := generator.Generate(ctx, gen, source)
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
	generateCmdFlags.StringVar(&outputPath, "output", "./out", "Output path for generated code")
	generateCmdFlags.BoolVar(&createDir, "create-dir", true, "Create output directory if it does not exist")
	generateCmdFlags.BoolVarP(&confirmYes, "yes", "y", false, "Skip the confirmation prompt and start generation immediately")
	generateCmdFlags.BoolVar(&useCache, "cache", true, "Use endpoint cache when generating; automatically creates it if missing")
	generateCmdFlags.BoolVar(&goMod, "go-mod", false, "Generate go.mod for Go output")
	generateCmdFlags.StringSliceVar(&includeEndpoints, "include-endpoints", []string{}, "Include only specified endpoints (comma-separated response keys)")
	generateCmdFlags.StringSliceVar(&excludeEndpoints, "exclude-endpoints", []string{}, "Exclude specified endpoints (comma-separated response keys)")
}
