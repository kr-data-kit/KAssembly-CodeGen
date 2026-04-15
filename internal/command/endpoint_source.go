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
	for result := range results {
		if result.Error != nil {
			return nil, result.Error
		}
		if result.Endpoint == nil {
			continue
		}
		services = append(services, result.Endpoint)
	}

	return services, nil
}
