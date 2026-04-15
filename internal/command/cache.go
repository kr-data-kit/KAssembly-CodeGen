package command

import (
	"context"
	"kassemblycodegen/internal/endpoint"
	"log/slog"

	"github.com/spf13/cobra"
)

var cacheCmd = &cobra.Command{
	Use:   "cache",
	Short: "Cache OpenAssembly service metadata",
	Long: `Fetch and cache OpenAssembly service summaries, specifications, and metadata.

The cache is written to kasm.cache in the current working directory.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithCancel(cmd.Context())
		defer cancel()

		slog.Info("Starting cache command", "file", endpoint.CacheFileName)
		cache, err := endpoint.UpdateCache(ctx, endpoint.CacheFileName, nil, nil)
		if err != nil {
			return err
		}

		slog.Info("Cache file written", "file", endpoint.CacheFileName, "services", len(cache.Services))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(cacheCmd)
}
