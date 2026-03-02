package command

import (
	"context"
	"fmt"
	"log/slog"
	"openassemblybinder/internal/service"

	"github.com/spf13/cobra"
)

var method string

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List OpenAssembly API services",
	Long: `List all available services from the OpenAssembly API.

Supported methods:
  simple   - Show basic service information (default)
  detailed - Show detailed service information

Example:
  openassemblybinder list --method simple
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// TODO: add ctx config
		ctx := context.Background()

		// Check if feature is available
		slog.Warn("This feature is still under development")

		data, err := service.FetchSummary(ctx)
		if err != nil {
			return err
		}

		switch method {
		case "detailed":
			slog.Info("Listing services (detailed mode)")
			return listDetailed(ctx, data)
		case "simple":
			slog.Info("Listing services (simple mode)")
			return listSimple(ctx, data)
		default:
			return fmt.Errorf("unknown method: %s", method)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	listCmdFlags := listCmd.Flags()
	listCmdFlags.StringVar(&method, "method", "simple", "Method to list services: simple or detailed")
}

func listSimple(
	ctx context.Context,
	data []service.Summary,
) error {
	// TODO : add list implementation
	return nil
}

func listDetailed(
	ctx context.Context,
	data []service.Summary,
) error {
	// TODO : add detailed list implementation
	return nil
}
