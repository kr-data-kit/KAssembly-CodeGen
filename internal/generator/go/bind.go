package gogen

import (
	"fmt"
	"kassemblycodegen/internal/service"
	"path"
	"strings"
)

type BindTemplateData struct {
	GlobalTemplateData
	Service *service.Service
}

func ExecuteBindTemplate(dir string, data BindTemplateData) error {
	bind := path.Join(dir, createFileName(data.Service.StructName))
	err := ExecuteTemplate(BindTemplate, bind, data)
	if err != nil {
		return fmt.Errorf("failed to execute bind template: %w", err)
	}
	return nil
}

func createFileName(name string) string {
	return "bind-" + strings.ToLower(name) + ".go"
}
