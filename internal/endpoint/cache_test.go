package endpoint

import (
	"context"
	"path/filepath"
	"testing"
)

func TestSaveAndLoadCache(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, CacheFileName)

	cache := &Cache{
		Services: []*Endpoint{
			{ID: "B", ResponseKey: "B", StructName: "B"},
			{ID: "A", ResponseKey: "A", StructName: "A"},
		},
	}

	if err := SaveCache(filePath, cache); err != nil {
		t.Fatalf("SaveCache() error = %v", err)
	}

	loaded, err := LoadCache(filePath)
	if err != nil {
		t.Fatalf("LoadCache() error = %v", err)
	}

	if loaded.GeneratedAt.IsZero() {
		t.Fatal("expected generated timestamp to be set")
	}
	if got := len(loaded.Services); got != 2 {
		t.Fatalf("expected 2 services, got %d", got)
	}
	if got := loaded.Services[0].ID; got != "A" {
		t.Fatalf("expected services to be sorted by ID, got first ID %q", got)
	}
}

func TestMergeCachedServices(t *testing.T) {
	existing := []*Endpoint{
		{ID: "A", ResponseKey: "A", StructName: "OldA"},
		{ID: "B", ResponseKey: "B", StructName: "B"},
	}
	updated := []*Endpoint{
		{ID: "A", ResponseKey: "A", StructName: "NewA"},
		{ID: "C", ResponseKey: "C", StructName: "C"},
	}

	merged := mergeCachedServices(existing, updated)
	if got := len(merged); got != 3 {
		t.Fatalf("expected 3 services, got %d", got)
	}
	if got := merged[0].StructName; got != "NewA" {
		t.Fatalf("expected updated service to replace existing one, got %q", got)
	}
	if got := merged[2].ID; got != "C" {
		t.Fatalf("expected merged list to include new service, got last ID %q", got)
	}
}

func TestCacheResultsFiltersEndpoints(t *testing.T) {
	cache := &Cache{
		Services: []*Endpoint{
			{ID: "A", ResponseKey: "A", StructName: "A"},
			{ID: "B", ResponseKey: "B", StructName: "B"},
		},
	}

	results, err := cache.Results(context.Background(), []string{"B"}, nil)
	if err != nil {
		t.Fatalf("Results() error = %v", err)
	}

	var endpoints []*Endpoint
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
	if got := endpoints[0].ID; got != "B" {
		t.Fatalf("expected endpoint B, got %q", got)
	}
}

func TestCacheResultsCancelledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	cache := &Cache{
		Services: []*Endpoint{{ID: "A", ResponseKey: "A", StructName: "A"}},
	}

	results, err := cache.Results(ctx, nil, nil)
	if err != nil {
		t.Fatalf("Results() error = %v", err)
	}

	result, ok := <-results
	if !ok {
		t.Fatal("expected cancellation result before channel close")
	}
	if result.Error == nil {
		t.Fatal("expected cancellation error")
	}
}
