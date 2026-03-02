package pygen

import (
	"embed"
	_ "embed"
	"fmt"
	"io/fs"
	"strings"
	"text/template"
)

type TemplateName string

//go:embed templates/*.tmpl
var templateFS embed.FS

var (
	TmplSet *template.Template
)

const (
	ClientTemplate        TemplateName = "client.tmpl"
	ModelsTemplate        TemplateName = "models.tmpl"
	StatusTemplate        TemplateName = "status.tmpl"
	ExceptionsTemplate    TemplateName = "exceptions.tmpl"
	InitTemplate          TemplateName = "init.tmpl"
	EndpointTemplate      TemplateName = "endpoint.tmpl"
	EndpointInitTemplate  TemplateName = "endpoint-init.tmpl"
	PyprojectTemplate     TemplateName = "pyproject.toml.tmpl"
	PythonVersionTemplate TemplateName = "python-version.tmpl"
)

func init() {
	var funcMap = template.FuncMap{
		"ToLowerCase": ToLowerCase,
		"RemoveLF":    RemoveLF,
	}

	subFS, err := fs.Sub(templateFS, "templates")
	if err != nil {
		panic(fmt.Errorf("error while generating subFS: %w", err))
	}

	TmplSet = template.Must(template.New("base").
		Funcs(funcMap).
		ParseFS(subFS, "*.tmpl"))
}

func GetTemplate(name TemplateName) (*template.Template, error) {
	// Lookup은 템플릿이 없으면 nil을 반환합니다.
	t := TmplSet.Lookup(string(name))
	if t == nil {
		return nil, fmt.Errorf("template %s not found in TmplSet", name)
	}
	return t, nil
}

// region Helper Functions

func ToLowerCase(input string) string {
	return strings.ToLower(input)
}

func RemoveLF(s string) string {
	s = strings.ReplaceAll(s, "\r", "")
	return strings.ReplaceAll(s, "\n", "")
}

// endregion
