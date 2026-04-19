package command

import (
	"bytes"
	"strings"
	"testing"
)

func TestReadGenerationConfirmationYes(t *testing.T) {
	var output bytes.Buffer

	confirmed, err := readGenerationConfirmation(strings.NewReader("y\n"), &output)
	if err != nil {
		t.Fatalf("readGenerationConfirmation() error = %v", err)
	}
	if !confirmed {
		t.Fatal("expected confirmation to be true")
	}
	if !strings.Contains(output.String(), "Proceed with code generation? [y/N]:") {
		t.Fatalf("expected prompt output, got %q", output.String())
	}
}

func TestReadGenerationConfirmationNo(t *testing.T) {
	var output bytes.Buffer

	confirmed, err := readGenerationConfirmation(strings.NewReader("n\n"), &output)
	if err != nil {
		t.Fatalf("readGenerationConfirmation() error = %v", err)
	}
	if confirmed {
		t.Fatal("expected confirmation to be false")
	}
}

func TestReadGenerationConfirmationEOF(t *testing.T) {
	var output bytes.Buffer

	_, err := readGenerationConfirmation(strings.NewReader(""), &output)
	if err == nil {
		t.Fatal("expected error for EOF input")
	}
	if !strings.Contains(err.Error(), "non-interactive stdin") {
		t.Fatalf("unexpected error: %v", err)
	}
}
