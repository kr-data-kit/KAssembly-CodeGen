package command

import (
	"context"
	"fmt"
	"log/slog"
	"sort"

	"github.com/spf13/cobra"
)

var listUseCache bool

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List available OpenAssembly APIs",
	Long: `Fetch and display the available OpenAssembly API list.

The command can use the local kasm.cache file or refresh it automatically when needed.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithCancel(cmd.Context())
		defer cancel()

		slog.Info("Starting API list command", "cache", listUseCache)

		source := newEndpointSource(listUseCache, nil, nil)
		services, err := collectEndpointsFromSource(ctx, source)
		if err != nil {
			return fmt.Errorf("failed to list APIs: %w", err)
		}

		sort.Slice(services, func(i, j int) bool {
			return services[i].ResponseKey < services[j].ResponseKey
		})

		fmt.Println("--------------------------------------------------------------------------------")
		fmt.Printf("Available APIs (%d)\n", len(services))
		fmt.Println("--------------------------------------------------------------------------------")
		fmt.Printf("%-28s %-56s\n", "ResponseKey", "Title")
		fmt.Println("--------------------------------------------------------------------------------")

		for _, service := range services {
			responseKey := truncateText(service.ResponseKey, 28)
			title := truncateText(service.Title, 56)
			fmt.Printf("%-28s %-56s\n", responseKey, title)
		}

		fmt.Println("--------------------------------------------------------------------------------")
		slog.Info("API list completed", "services", len(services))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().BoolVar(&listUseCache, "cache", true, "Use endpoint cache when listing APIs; automatically creates it if missing")
}

func truncateText(value string, limit int) string {
	runes := []rune(value)
	if len(runes) <= limit {
		return value
	}
	if limit <= 1 {
		return string(runes[:limit])
	}
	return string(runes[:limit-1]) + "..."
}
