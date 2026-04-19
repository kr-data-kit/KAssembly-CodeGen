package command

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

var listUseCache bool
var listExtraFields string

type listExtraConfig struct {
	showID          bool
	showRequestArgs bool
	showResultArgs  bool
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List available OpenAssembly APIs",
	Long: `Fetch and display the available OpenAssembly API list.

The command can use the local kasm.cache file or refresh it automatically when needed.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithCancel(cmd.Context())
		defer cancel()

		slog.Info("Starting API list command", "cache", listUseCache)

		extraConfig, err := parseListExtraFields(listExtraFields)
		if err != nil {
			return err
		}

		source := newEndpointSource(listUseCache, nil, nil)
		services, err := collectEndpointsFromSource(ctx, source)
		if err != nil {
			return fmt.Errorf("failed to list APIs: %w", err)
		}

		sort.Slice(services, func(i, j int) bool {
			return services[i].ResponseKey < services[j].ResponseKey
		})

		headers := []string{"ResponseKey", "Title"}
		separators := []string{"-----------", "-----"}
		if extraConfig.showID {
			headers = append(headers, "ID")
			separators = append(separators, "--")
		}
		if extraConfig.showRequestArgs {
			headers = append(headers, "RequestArgs")
			separators = append(separators, "-----------")
		}
		if extraConfig.showResultArgs {
			headers = append(headers, "ResultArgs")
			separators = append(separators, "----------")
		}

		writer := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(writer, strings.Join(headers, "\t"))
		fmt.Fprintln(writer, strings.Join(separators, "\t"))

		for _, service := range services {
			row := []string{
				service.ResponseKey,
				truncateText(service.Title, 40),
			}
			if extraConfig.showID {
				row = append(row, truncateText(service.ID, 20))
			}
			if extraConfig.showRequestArgs {
				row = append(row, truncateText(formatRequestArgs(service.Params), 90))
			}
			if extraConfig.showResultArgs {
				row = append(row, truncateText(formatResultArgs(service.Cols), 90))
			}

			fmt.Fprintln(writer, strings.Join(row, "\t"))
		}

		if err := writer.Flush(); err != nil {
			return fmt.Errorf("failed to flush API list output: %w", err)
		}
		
		slog.Info("API list completed", "services", len(services))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().BoolVar(&listUseCache, "cache", true, "Use endpoint cache when listing APIs; automatically creates it if missing")
	listCmd.Flags().StringVar(&listExtraFields, "extra", "", "Additional columns: id, request-args, result-args (comma-separated)")
}

func parseListExtraFields(raw string) (listExtraConfig, error) {
	config := listExtraConfig{}
	if strings.TrimSpace(raw) == "" {
		return config, nil
	}

	for _, field := range strings.Split(raw, ",") {
		normalized := strings.TrimSpace(strings.ToLower(field))
		if normalized == "" {
			continue
		}

		switch normalized {
		case "id":
			config.showID = true
		case "request-args":
			config.showRequestArgs = true
		case "result-args":
			config.showResultArgs = true
		default:
			return listExtraConfig{}, fmt.Errorf("invalid --extra value %q (allowed: id, request-args, result-args)", normalized)
		}
	}

	return config, nil
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
