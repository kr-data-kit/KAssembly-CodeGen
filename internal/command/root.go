package command

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "openassemblybinder",
	Short: "OpenAssembly API Client Code Generator",
	Long: `OpenAssembly-Binder generates client code for the National Assembly of Korea's Open API.

It fetches the latest API specifications and generates type-safe client code
for Go and Python languages.

For more information, visit: https://github.com/rethinking21/OpenAssembly-Binder`,
	Version: "0.1.0",
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}
