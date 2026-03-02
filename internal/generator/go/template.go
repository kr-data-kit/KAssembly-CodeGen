package gogen

import (
	"fmt"
	"os"
)

type GlobalTemplateData struct {
	PackageName   string
	ClientName    string // 필요 없어짐
	RepositoryURL string
}

func ExecuteTemplate(templateName TemplateName, file string, data any) error {
	tmpl, err := GetTemplate(templateName)
	if err != nil {
		return fmt.Errorf("failed to get template %s: %w", templateName, err)
	}

	f, err := os.Create(file)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", file, err)
	}
	defer f.Close()

	err = tmpl.Execute(f, data)
	if err != nil {
		return fmt.Errorf("failed to execute template %s: %w", templateName, err)
	}
	return nil
}
