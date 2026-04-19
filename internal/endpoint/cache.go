package endpoint

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"sort"
	"time"
)

const CacheFileName = "kasm.cache"

type Cache struct {
	GeneratedAt time.Time   `json:"generatedAt"`
	Services    []*Endpoint `json:"services"`
}

func (c *Cache) Results(ctx context.Context, includeList, excludeList []string) (<-chan *GenerateResult, error) {
	if c == nil {
		return nil, fmt.Errorf("cache is required")
	}

	resultChan := make(chan *GenerateResult)
	go func() {
		defer close(resultChan)

		for _, service := range c.Services {
			if ctx.Err() != nil {
				resultChan <- &GenerateResult{Error: fmt.Errorf("cache stream cancelled: %w", ctx.Err())}
				return
			}
			if service == nil {
				continue
			}
			if !hasRequestedMatch(service.ID, service.ResponseKey, includeList, excludeList) {
				continue
			}

			resultChan <- &GenerateResult{Endpoint: service}
		}
	}()

	return resultChan, nil
}

func LoadCache(filePath string) (*Cache, error) {
	slog.Debug("Loading cache file", "file", filePath)

	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("%w: %s", os.ErrNotExist, filePath)
		}
		return nil, fmt.Errorf("failed to read cache file %s: %w", filePath, err)
	}

	var cache Cache
	if err := json.Unmarshal(data, &cache); err != nil {
		return nil, fmt.Errorf("failed to decode cache file %s: %w", filePath, err)
	}

	sortEndpoints(cache.Services)
	slog.Debug("Loaded cache file", "file", filePath, "services", len(cache.Services))
	return &cache, nil
}

func SaveCache(filePath string, cache *Cache) error {
	if cache == nil {
		return fmt.Errorf("cache is required")
	}

	slog.Debug("Saving cache file", "file", filePath, "services", len(cache.Services))
	cache.GeneratedAt = time.Now().UTC()
	sortEndpoints(cache.Services)

	data, err := json.MarshalIndent(cache, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to encode cache data: %w", err)
	}

	if err := os.WriteFile(filePath, data, 0o644); err != nil {
		return fmt.Errorf("failed to write cache file %s: %w", filePath, err)
	}

	slog.Info("Cache file saved", "file", filePath, "services", len(cache.Services))
	return nil
}

func UpdateCache(ctx context.Context, filePath string) (*Cache, error) {
	slog.Info("Updating cache", "file", filePath)

	currentCache, err := LoadCache(filePath)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return nil, err
		}
		slog.Info("Cache file not found, creating a new one", "file", filePath)
		currentCache = &Cache{}
	} else {
		slog.Debug("Merging existing cache entries", "file", filePath, "services", len(currentCache.Services))
	}

	services, err := collectEndpoints(ctx)
	if err != nil {
		return nil, err
	}
	slog.Info("Fetched live endpoints for cache update", "count", len(services))

	currentCache.Services = mergeCachedServices(currentCache.Services, services)
	if err := SaveCache(filePath, currentCache); err != nil {
		return nil, err
	}

	slog.Info("Cache update complete", "file", filePath, "services", len(currentCache.Services))
	return currentCache, nil
}

func collectEndpoints(ctx context.Context) ([]*Endpoint, error) {
	slog.Debug("Collecting endpoints for cache update")

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	results, err := GenerateEndpoints(ctx, nil, nil)
	if err != nil {
		return nil, err
	}

	services := make([]*Endpoint, 0)
	for result := range results {
		if result.Error != nil {
			return nil, result.Error
		}
		if result.Endpoint == nil {
			continue
		}
		slog.Debug("Collected endpoint", "id", result.Endpoint.ID, "responseKey", result.Endpoint.ResponseKey, "structName", result.Endpoint.StructName)
		services = append(services, result.Endpoint)
	}

	slog.Info("Collected endpoints", "count", len(services))
	return services, nil
}

func mergeCachedServices(existing, updated []*Endpoint) []*Endpoint {
	serviceMap := make(map[string]*Endpoint, len(existing)+len(updated))
	order := make([]string, 0, len(existing)+len(updated))

	addService := func(service *Endpoint) {
		if service == nil {
			return
		}
		key := cacheKey(service)
		if _, exists := serviceMap[key]; !exists {
			order = append(order, key)
		}
		serviceMap[key] = service
	}

	for _, service := range existing {
		addService(service)
	}
	for _, service := range updated {
		addService(service)
	}

	merged := make([]*Endpoint, 0, len(serviceMap))
	for _, key := range order {
		merged = append(merged, serviceMap[key])
	}

	sortEndpoints(merged)
	return merged
}

func sortEndpoints(services []*Endpoint) {
	sort.Slice(services, func(i, j int) bool {
		if services[i].ID == services[j].ID {
			return services[i].ResponseKey < services[j].ResponseKey
		}
		return services[i].ID < services[j].ID
	})
}

func cacheKey(service *Endpoint) string {
	if service == nil {
		return ""
	}
	return fmt.Sprintf("%s|%s", service.ID, service.ResponseKey)
}
