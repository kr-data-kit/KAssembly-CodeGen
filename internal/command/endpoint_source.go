package command

import (
	"context"
	"fmt"
	"kassemblycodegen/internal/endpoint"
	"kassemblycodegen/internal/generator"
)

func newEndpointSource(useCache bool, includeList, excludeList []string) generator.EndpointSource {
	if useCache {
		return newCachedEndpointSource(includeList, excludeList, true)
	}

	return generator.LiveEndpointSource{
		IncludeList: includeList,
		ExcludeList: excludeList,
	}
}

func newCachedEndpointSource(includeList, excludeList []string, autoUpdate bool) generator.EndpointSource {
	return generator.CachedEndpointSource{
		CacheFile:   endpoint.CacheFileName,
		IncludeList: includeList,
		ExcludeList: excludeList,
		AutoUpdate:  autoUpdate,
	}
}

func collectEndpointsFromSource(ctx context.Context, source generator.EndpointSource) ([]*endpoint.Endpoint, error) {
	results, err := source.Generate(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare endpoints: %w", err)
	}

	services := make([]*endpoint.Endpoint, 0)
	var firstErr error
	for result := range results {
		if result == nil {
			continue
		}
		if result.Error != nil {
			if firstErr == nil {
				firstErr = result.Error
			}
			continue
		}
		if result.Endpoint == nil {
			continue
		}
		services = append(services, result.Endpoint)
	}

	if firstErr != nil {
		return nil, firstErr
	}

	return services, nil
}
