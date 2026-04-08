package gogen

import (
	"embed"
	"fmt"
	"io/fs"
	"kassemblycodegen/internal/generator"
	"text/template"
)

//go:embed all:templates
var embedFS embed.FS

var (
	ClientTemplates *template.Template // client
)

func init() {
	subFS, err := fs.Sub(embedFS, "templates")
	if err != nil {
		panic(fmt.Errorf("error while generating subFS: %w", err))
	}

	ClientTemplates = generator.CreateTemplateSet(subFS, "client")
}
