package command

import (
	"context"
	"fmt"
	"kassemblycodegen/internal/endpoint"
	"log/slog"
	"sort"
	"strings"

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
	lines := []string{
		fmt.Sprintf("Response Key: %s", service.ResponseKey),
		fmt.Sprintf("Title: %s", service.Title),
		fmt.Sprintf("ID: %s", service.ID),
		fmt.Sprintf("URL: %s", service.URL),
		fmt.Sprintf("Description: %s", service.Description),
		fmt.Sprintf("Request Args: %s", formatRequestArgs(service.Params)),
		fmt.Sprintf("Result Args: %s", formatResultArgs(service.Cols)),
		fmt.Sprintf("Provides API: %t", service.ProvidesAPI),
		fmt.Sprintf("Provides Data: %t", service.ProvidesData),
		fmt.Sprintf("Commercial Use Allowed: %t", service.CommercialUseAllowed),
		fmt.Sprintf("Attribution Required: %t", service.AttributionRequired),
	}

	maxWidth := len("API Information")
	for _, line := range lines {
		if len(line) > maxWidth {
			maxWidth = len(line)
		}
	}

	fmt.Println("--------------------------------------------------------------------------------")
	fmt.Println("API Information")
	fmt.Println("--------------------------------------------------------------------------------")
	for _, line := range lines {
		fmt.Println(line)
	}
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
