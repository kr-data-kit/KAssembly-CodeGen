package command

import "testing"

func TestValidateGoPackageName(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{name: "valid", value: "openassemblyclient", wantErr: false},
		{name: "underscore", value: "_client", wantErr: false},
		{name: "empty", value: "", wantErr: true},
		{name: "hyphen", value: "my-client", wantErr: true},
		{name: "keyword", value: "map", wantErr: true},
		{name: "digit start", value: "1client", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateGoPackageName(tt.value)
			if tt.wantErr && err == nil {
				t.Fatalf("expected error for %q", tt.value)
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("unexpected error for %q: %v", tt.value, err)
			}
		})
	}
}

func TestValidatePythonPackageName(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{name: "valid", value: "openassemblyclient", wantErr: false},
		{name: "underscore", value: "_client", wantErr: false},
		{name: "empty", value: "", wantErr: true},
		{name: "hyphen", value: "my-client", wantErr: true},
		{name: "keyword", value: "class", wantErr: true},
		{name: "digit start", value: "1client", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validatePythonPackageName(tt.value)
			if tt.wantErr && err == nil {
				t.Fatalf("expected error for %q", tt.value)
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("unexpected error for %q: %v", tt.value, err)
			}
		})
	}
}

func TestValidatePackageNameDispatch(t *testing.T) {
	if err := validatePackageName("go", "openassemblyclient"); err != nil {
		t.Fatalf("unexpected error for go: %v", err)
	}
	if err := validatePackageName("python", "openassemblyclient"); err != nil {
		t.Fatalf("unexpected error for python: %v", err)
	}
	if err := validatePackageName("ruby", "openassemblyclient"); err == nil {
		t.Fatal("expected unsupported language error")
	}
}
