package generator

import (
	"context"
	"errors"
	"fmt"
	"kassemblycodegen/internal/endpoint"
	"os"
)

type LiveEndpointSource struct {
	IncludeList []string
	ExcludeList []string
}

func (s LiveEndpointSource) Generate(ctx context.Context) (<-chan *endpoint.GenerateResult, error) {
	return endpoint.GenerateEndpoints(ctx, s.IncludeList, s.ExcludeList)
}

type CachedEndpointSource struct {
	CacheFile   string
	IncludeList []string
	ExcludeList []string
	AutoUpdate  bool
}

func (s CachedEndpointSource) Generate(ctx context.Context) (<-chan *endpoint.GenerateResult, error) {
	cache, err := endpoint.LoadCache(s.CacheFile)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return nil, err
		}
		if !s.AutoUpdate {
			return nil, fmt.Errorf("cache file not found: %s", s.CacheFile)
		}

		cache, err = endpoint.UpdateCache(ctx, s.CacheFile, s.IncludeList, s.ExcludeList)
		if err != nil {
			return nil, err
		}
	}

	return cache.Results(ctx, s.IncludeList, s.ExcludeList)
}
