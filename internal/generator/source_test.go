package generator

import (
	"context"
	"path/filepath"
	"testing"

	"kassemblycodegen/internal/endpoint"
)

func TestCachedEndpointSourceGenerate(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, endpoint.CacheFileName)

	cache := &endpoint.Cache{
		Services: []*endpoint.Endpoint{
			{ID: "A", ResponseKey: "A", StructName: "A"},
			{ID: "B", ResponseKey: "B", StructName: "B"},
		},
	}
	if err := endpoint.SaveCache(filePath, cache); err != nil {
		t.Fatalf("SaveCache() error = %v", err)
	}

	source := CachedEndpointSource{
		CacheFile:   filePath,
		IncludeList: []string{"A"},
		ExcludeList: nil,
		AutoUpdate:  false,
	}

	results, err := source.Generate(context.Background())
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	var endpoints []*endpoint.Endpoint
	for result := range results {
		if result.Error != nil {
			t.Fatalf("unexpected result error: %v", result.Error)
		}
		if result.Endpoint != nil {
			endpoints = append(endpoints, result.Endpoint)
		}
	}

	if got := len(endpoints); got != 1 {
		t.Fatalf("expected 1 endpoint, got %d", got)
	}
	if got := endpoints[0].ID; got != "A" {
		t.Fatalf("expected endpoint A, got %q", got)
	}
}

func TestCachedEndpointSourceMissingFileWithoutAutoUpdate(t *testing.T) {
	source := CachedEndpointSource{
		CacheFile:  filepath.Join(t.TempDir(), endpoint.CacheFileName),
		AutoUpdate: false,
	}

	_, err := source.Generate(context.Background())
	if err == nil {
		t.Fatal("expected error when cache file is missing")
	}
}
