package generator

import (
	"context"
	"fmt"
	"kassemblycodegen/internal/endpoint"
	"log/slog"
)

const (
	RepositoryURL = "https://github.com/kr-data-kit/KAssembly-CodeGen"
)

type Generator interface {
	SetGlobalConfig() error // thinking

	GenStatic() error
	GenEndpoint(ctx context.Context, enp *endpoint.Endpoint) error
	GenFinal() error
}

type GeneralConfig struct {
	PackageName     string
	OutPath         string
	CreateDir       bool
	IncludeEndpoint []string
	ExcludeEndpoint []string
}

func Generate(ctx context.Context, gen Generator, includeList, excludeList []string) error {
	err := gen.SetGlobalConfig()
	if err != nil {
		return fmt.Errorf("failed to set global config: %w", err)
	}

	err = gen.GenStatic()
	if err != nil {
		return fmt.Errorf("failed to generate static files: %w", err)
	}

	endpoints, err := endpoint.GenerateEndpoints(ctx, includeList, excludeList)
	if err != nil {
		return fmt.Errorf("failed to generate endpoints: %w", err)
	}

	running := true

	for running {
		select {
		case <-ctx.Done():
			running = false
			continue
		case result, ok := <-endpoints:
			if !ok {
				// channel closed, all services processed
				running = false
				continue
			}
			if result.Error != nil {
				slog.Error("Error generating endpoint", "error", result.Error)
				continue
			}

			endpoint := result.Endpoint

			// TODO: filter?
			err = gen.GenEndpoint(ctx, endpoint)
			if err != nil {
				slog.Error("Error generating endpoint", "endpoint", endpoint.StructName, "error", err)
				continue
			}
		}
	}

	err = gen.GenFinal()
	if err != nil {
		return fmt.Errorf("failed to generate final files: %w", err)
	}
	return nil
}
