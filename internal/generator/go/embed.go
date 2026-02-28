package gogen

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
	BindTemplate   TemplateName = "bind.tmpl"
	ClientTemplate TemplateName = "client.tmpl"
	CommonTemplate TemplateName = "common.tmpl"
	StatusTemplate TemplateName = "status.tmpl"
	GoModTemplate  TemplateName = "go.mod.tmpl"
)

func init() {
	var funcMap = template.FuncMap{
		"KebabToPascalCase":       KebabToPascalCase,
		"ScreamSnakeToPascalCase": ScreamSnakeToPascalCase,
		"RemoveLF":                RemoveLF,
		"ToLowerCase":             ToLowerCase,
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

func KebabToPascalCase(input string) string {
	parts := strings.Split(input, "-")

	var b strings.Builder
	for _, p := range parts {
		if p == "" {
			continue
		}
		b.WriteString(strings.ToUpper(p[:1]))
		if len(p) > 1 {
			b.WriteString(p[1:])
		}
	}
	return b.String()
}

func ScreamSnakeToPascalCase(input string) string {
	parts := strings.Split(strings.ToLower(input), "_")

	var b strings.Builder
	for _, p := range parts {
		if p == "" {
			continue
		}
		b.WriteString(strings.ToUpper(p[:1]))
		if len(p) > 1 {
			b.WriteString(p[1:])
		}
	}
	return b.String()
}

func ToLowerCase(input string) string {
	return strings.ToLower(input)
}

func RemoveLF(s string) string {
	s = strings.ReplaceAll(s, "\r", "")
	return strings.ReplaceAll(s, "\n", "")
}

// endregion
