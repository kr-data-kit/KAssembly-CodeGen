package endpoint

import (
	"context"
	"testing"
)

func TestGetStructName(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{name: "lowercase", in: "abc", want: "Abc"},
		{name: "already mixed", in: "ALLBILL", want: "ALLBILL"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getStructName(tt.in); got != tt.want {
				t.Fatalf("getStructName(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestServiceTypeHelpers(t *testing.T) {
	apiProvides, apiInfSeq, dataProvides, dataInfSeq := parseServiceTypes("A-2,S-1")
	if !apiProvides || apiInfSeq != "2" {
		t.Fatalf("unexpected API parse result: provides=%v seq=%q", apiProvides, apiInfSeq)
	}
	if !dataProvides || dataInfSeq != "1" {
		t.Fatalf("unexpected data parse result: provides=%v seq=%q", dataProvides, dataInfSeq)
	}

	if !shouldIncludeEndpoint("A", nil, nil) {
		t.Fatal("expected include when both lists are empty")
	}
	if shouldIncludeEndpoint("A", []string{"A"}, []string{"A"}) {
		t.Fatal("expected exclude list to win")
	}
	if !hasRequestedMatch("A", "B", []string{"B"}, nil) {
		t.Fatal("expected response key match to be accepted")
	}
	if !isExcluded("A", []string{"A"}) {
		t.Fatal("expected ID to be excluded")
	}
}

func TestApplyQueryMetadata(t *testing.T) {
	enp := &Endpoint{}
	applyQueryMetadata(enp, "상업용 금지 / 출처표시 필요")

	if enp.CCL != "상업용 금지 / 출처표시 필요" {
		t.Fatalf("unexpected CCL: %q", enp.CCL)
	}
	if !enp.CommercialUseAllowed {
		t.Fatal("expected commercial use to be allowed flag to be true when license contains it")
	}
	if !enp.AttributionRequired {
		t.Fatal("expected attribution required flag to be true when license contains it")
	}
}

func TestNewEndpointFromSummary(t *testing.T) {
	item := Summary{
		ID:          "TESTID",
		Title:       "Test Title",
		Description: "Test Description",
	}
	spec := &ServiceSpec{
		Endpoint:    "TESTENDPOINT",
		ResponseKey: "testresponse",
		Variables:   []Variable{{ID: "A"}},
		Columns:     []Column{{ID: "B"}},
	}

	enp := newEndpointFromSummary(item, spec)
	if enp.ID != item.ID {
		t.Fatalf("unexpected ID: %q", enp.ID)
	}
	if enp.URL == "" {
		t.Fatal("expected URL to be populated")
	}
	if enp.StructName != "Testresponse" {
		t.Fatalf("unexpected struct name: %q", enp.StructName)
	}
	if len(enp.AlterStructNames) != 1 || enp.AlterStructNames[0] != item.ID {
		t.Fatalf("unexpected alter struct names: %#v", enp.AlterStructNames)
	}
	if enp.Endpoint != spec.Endpoint || enp.ResponseKey != spec.ResponseKey {
		t.Fatalf("unexpected endpoint mapping: %+v", enp)
	}
	if len(enp.Params) != 1 || len(enp.Cols) != 1 {
		t.Fatalf("unexpected params/cols: %+v", enp)
	}
}

func TestNewVersionedEndpoint(t *testing.T) {
	base := &Endpoint{
		ID:                   ALLBILL,
		Title:                "Base Title",
		Description:          "Base Description",
		CCL:                  "license",
		CommercialUseAllowed: true,
		AttributionRequired:  true,
	}
	spec := &ServiceSpec{
		Endpoint:    "NEW-ENDPOINT",
		ResponseKey: "newresponse",
		Variables:   []Variable{{ID: "A"}},
		Columns:     []Column{{ID: "B"}},
	}

	enp := newVersionedEndpoint(base, spec, BILLRCPV, "BASE", " (Version 1)", "https://example.com", "2", "1")
	if enp.ID != BILLRCPV {
		t.Fatalf("unexpected ID: %q", enp.ID)
	}
	if enp.Title != "Base Title (Version 1)" {
		t.Fatalf("unexpected title: %q", enp.Title)
	}
	if enp.Description != "Base Description - Updated Version" {
		t.Fatalf("unexpected description: %q", enp.Description)
	}
	if enp.Endpoint != spec.Endpoint || enp.ResponseKey != spec.ResponseKey {
		t.Fatalf("unexpected endpoint mapping: %+v", enp)
	}
	if enp.APIInfSeq != "2" || enp.DataInfSeq != "1" {
		t.Fatalf("unexpected sequences: %+v", enp)
	}
	if !enp.CommercialUseAllowed || !enp.AttributionRequired {
		t.Fatal("expected license flags to be copied from base endpoint")
	}
}

func TestCheckAdditionalIDsValidation(t *testing.T) {
	tests := []string{"", "unsafe/id", "unsafe?id", "unsafe#id"}
	for _, value := range tests {
		if _, err := CheckAdditionalIDs(context.Background(), value); err == nil {
			t.Fatalf("expected validation error for %q", value)
		}
	}
}
