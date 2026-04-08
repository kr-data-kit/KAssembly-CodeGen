package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type TemplateEngine struct {
	template *template.Template
	root     string
}

func NewTemplateEngine(template *template.Template, root string) *TemplateEngine {
	// TODO : ensure root path exist
	return &TemplateEngine{
		template: template,
		root:     root,
	}
}

func (e *TemplateEngine) Execute(name string, filePath string, data any) error {
	tmpl := e.template.Lookup(name)
	if tmpl == nil {
		return fmt.Errorf("template %s not found in TmplSet", name)
	}

	fullPath := filepath.Join(e.root, filePath)
	if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
		return fmt.Errorf("failed to create directory for %s: %w", fullPath, err)
	}

	f, err := os.Create(fullPath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", fullPath, err)
	}
	defer f.Close()

	err = tmpl.Execute(f, data)
	if err != nil {
		return fmt.Errorf("failed to execute template %s: %w", name, err)
	}
	return nil
}

func (e *TemplateEngine) GenerateFile(name string, data any) error {
	fileName := strings.TrimSuffix(name, ".tmpl")
	return e.Execute(name, fileName, data)
}

func (e *TemplateEngine) GenerateAll(names []string, data any) error {
	for _, name := range names {
		err := e.GenerateFile(name, data)
		if err != nil {
			return err
		}
	}
	return nil
}
