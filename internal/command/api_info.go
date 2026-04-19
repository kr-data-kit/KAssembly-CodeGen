package command

import (
	"context"
	"fmt"
	"kassemblycodegen/internal/endpoint"
	"log/slog"
	"os"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

var apiInfoUseCache bool

var apiInfoCmd = &cobra.Command{
	Use:   "api-info <ResponseKey>",
	Short: "Show detailed OpenAssembly API information",
	Long: `Display a detailed summary for a single OpenAssembly API.

This command reads from kasm.cache only and will not refresh cache data automatically.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithCancel(cmd.Context())
		defer cancel()

		responseKey := args[0]
		slog.Info("Starting API info command", "responseKey", responseKey, "cache", apiInfoUseCache)

		if !apiInfoUseCache {
			return fmt.Errorf("api-info requires cache mode; enable --cache")
		}

		source := newCachedEndpointSource(nil, nil, false)
		services, err := collectEndpointsFromSource(ctx, source)
		if err != nil {
			return fmt.Errorf("failed to load API info: %w", err)
		}

		matched := filterEndpointsByResponseKey(services, responseKey)
		if len(matched) == 0 {
			return fmt.Errorf("api not found in cache: %s", responseKey)
		}

		sort.Slice(matched, func(i, j int) bool {
			if matched[i].ID == matched[j].ID {
				return matched[i].ResponseKey < matched[j].ResponseKey
			}
			return matched[i].ID < matched[j].ID
		})

		for index, service := range matched {
			if index > 0 {
				fmt.Println()
			}
			printAPIInfoCard(service)
		}

		slog.Info("API info completed", "responseKey", responseKey, "matches", len(matched))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(apiInfoCmd)
	apiInfoCmd.Flags().BoolVar(&apiInfoUseCache, "cache", true, "Read API info from cache only")
}

func filterEndpointsByResponseKey(services []*endpoint.Endpoint, responseKey string) []*endpoint.Endpoint {
	matched := make([]*endpoint.Endpoint, 0)
	for _, service := range services {
		if service == nil {
			continue
		}
		if strings.EqualFold(service.ResponseKey, responseKey) {
			matched = append(matched, service)
		}
	}
	return matched
}

func printAPIInfoCard(service *endpoint.Endpoint) {
	rows := [][2]string{
		{"Response Key", service.ResponseKey},
		{"Title", service.Title},
		{"ID", service.ID},
		{"URL", service.URL},
		{"Description", service.Description},
		{"Request Args", formatRequestArgs(service.Params)},
		{"Result Args", formatResultArgs(service.Cols)},
		{"Provides API", fmt.Sprintf("%t", service.ProvidesAPI)},
		{"Provides Data", fmt.Sprintf("%t", service.ProvidesData)},
		{"Commercial Use Allowed", fmt.Sprintf("%t", service.CommercialUseAllowed)},
		{"Attribution Required", fmt.Sprintf("%t", service.AttributionRequired)},
	}

	fmt.Println("--------------------------------------------------------------------------------")
	fmt.Println("API Information")
	fmt.Println("--------------------------------------------------------------------------------")

	writer := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	for _, row := range rows {
		fmt.Fprintf(writer, "%s\t%s\n", row[0], row[1])
	}
	_ = writer.Flush()

	fmt.Println("--------------------------------------------------------------------------------")
}

func formatRequestArgs(args []endpoint.Variable) string {
	if len(args) == 0 {
		return "(none)"
	}

	parts := make([]string, 0, len(args))
	for _, arg := range args {
		parts = append(parts, fmt.Sprintf("%s(%s)", arg.ID, arg.Name))
	}

	return strings.Join(parts, ", ")
}

func formatResultArgs(args []endpoint.Column) string {
	if len(args) == 0 {
		return "(none)"
	}

	parts := make([]string, 0, len(args))
	for _, arg := range args {
		parts = append(parts, fmt.Sprintf("%s(%s)", arg.ID, arg.Name))
	}

	return strings.Join(parts, ", ")
}
