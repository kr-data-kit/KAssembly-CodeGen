package command

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"kassemblycodegen/internal/endpoint"
)

func TestFilterEndpointsByResponseKey(t *testing.T) {
	services := []*endpoint.Endpoint{
		{ID: "A", ResponseKey: "ALLBILL"},
		{ID: "B", ResponseKey: "OTHER"},
	}

	matched := filterEndpointsByResponseKey(services, "allbill")
	if got := len(matched); got != 1 {
		t.Fatalf("expected 1 match, got %d", got)
	}
	if got := matched[0].ID; got != "A" {
		t.Fatalf("expected match A, got %q", got)
	}
}

func TestPrintAPIInfoCard(t *testing.T) {
	service := &endpoint.Endpoint{
		ID:                   "TESTID",
		Title:                "Test Title",
		Description:          "Test Description",
		ResponseKey:          "TESTKEY",
		URL:                  "https://example.com",
		ProvidesAPI:          true,
		ProvidesData:         false,
		Params:               []endpoint.Variable{{ID: "AGE", Name: "나이"}},
		Cols:                 []endpoint.Column{{ID: "BILL_ID", Name: "의안ID"}},
		CommercialUseAllowed: false,
		AttributionRequired:  true,
	}

	output := captureStdout(t, func() {
		printAPIInfoCard(service)
	})

	for _, expected := range []string{"API Information", "Response Key", "TESTKEY", "Title", "Test Title", "Request Args", "AGE(나이)", "Result Args", "BILL_ID(의안ID)"} {
		if !strings.Contains(output, expected) {
			t.Fatalf("expected output to contain %q, got:\n%s", expected, output)
		}
	}
}

func TestTruncateTextKoreanSafe(t *testing.T) {
	input := "대한민국국회열린데이터"
	output := truncateText(input, 7)

	if output != "대한민국국회..." {
		t.Fatalf("unexpected truncate result: %q", output)
	}
	if !strings.HasSuffix(output, "...") {
		t.Fatalf("expected ellipsis suffix, got %q", output)
	}
}

func captureStdout(t *testing.T, fn func()) string {
	t.Helper()

	old := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe() error = %v", err)
	}
	os.Stdout = w

	outputCh := make(chan string, 1)
	go func() {
		var buf bytes.Buffer
		_, _ = io.Copy(&buf, r)
		outputCh <- buf.String()
	}()

	fn()

	_ = w.Close()
	os.Stdout = old

	return <-outputCh
}
