package command

import (
	"context"
	"errors"
	"strings"
	"testing"

	"kassemblycodegen/internal/endpoint"
	"kassemblycodegen/internal/generator"
)

type fakeEndpointSource struct {
	results []*endpoint.GenerateResult
	err     error
}

func (f fakeEndpointSource) Generate(ctx context.Context) (<-chan *endpoint.GenerateResult, error) {
	if f.err != nil {
		return nil, f.err
	}

	resultChan := make(chan *endpoint.GenerateResult, len(f.results))
	for _, result := range f.results {
		resultChan <- result
	}
	close(resultChan)
	return resultChan, nil
}

func TestNewEndpointSourceUsesCacheWhenEnabled(t *testing.T) {
	source := newEndpointSource(true, []string{"A"}, []string{"B"})

	cached, ok := source.(generator.CachedEndpointSource)
	if !ok {
		t.Fatalf("expected CachedEndpointSource, got %T", source)
	}
	if !cached.AutoUpdate {
		t.Fatal("expected cache auto update to be enabled")
	}
	if cached.CacheFile != endpoint.CacheFileName {
		t.Fatalf("unexpected cache file: %q", cached.CacheFile)
	}
}

func TestNewEndpointSourceUsesLiveWhenDisabled(t *testing.T) {
	source := newEndpointSource(false, []string{"A"}, []string{"B"})

	live, ok := source.(generator.LiveEndpointSource)
	if !ok {
		t.Fatalf("expected LiveEndpointSource, got %T", source)
	}
	if len(live.IncludeList) != 1 || live.IncludeList[0] != "A" {
		t.Fatalf("unexpected include list: %#v", live.IncludeList)
	}
}

func TestCollectEndpointsFromSource(t *testing.T) {
	source := fakeEndpointSource{
		results: []*endpoint.GenerateResult{
			{Endpoint: &endpoint.Endpoint{ID: "A", ResponseKey: "A"}},
			{Endpoint: &endpoint.Endpoint{ID: "B", ResponseKey: "B"}},
		},
	}

	services, err := collectEndpointsFromSource(context.Background(), source)
	if err != nil {
		t.Fatalf("collectEndpointsFromSource() error = %v", err)
	}
	if got := len(services); got != 2 {
		t.Fatalf("expected 2 services, got %d", got)
	}
}

func TestCollectEndpointsFromSourceReturnsSourceError(t *testing.T) {
	source := fakeEndpointSource{err: errors.New("boom")}

	_, err := collectEndpointsFromSource(context.Background(), source)
	if err == nil {
		t.Fatal("expected error from source")
	}
}

func TestParseListExtraFields(t *testing.T) {
	config, err := parseListExtraFields("id, request-args, result-args")
	if err != nil {
		t.Fatalf("parseListExtraFields() error = %v", err)
	}

	if !config.showID {
		t.Fatal("expected showID to be true")
	}
	if !config.showRequestArgs {
		t.Fatal("expected showRequestArgs to be true")
	}
	if !config.showResultArgs {
		t.Fatal("expected showResultArgs to be true")
	}
}

func TestParseListExtraFieldsRejectsInvalid(t *testing.T) {
	_, err := parseListExtraFields("id,unknown")
	if err == nil {
		t.Fatal("expected error for invalid value")
	}
	if !strings.Contains(err.Error(), "invalid --extra value") {
		t.Fatalf("unexpected error: %v", err)
	}
}
