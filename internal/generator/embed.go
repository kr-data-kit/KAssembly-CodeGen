package generator

import (
	"fmt"
	"io/fs"
	"strings"
	"text/template"
)

func CreateTemplateSet(subFS fs.FS, lang string) *template.Template {
	langFS, err := fs.Sub(subFS, lang)
	if err != nil {
		panic(fmt.Errorf("error while generating FS: %w", err))
	}

	ts := template.New(lang).Funcs(FuncMap)
	err = fs.WalkDir(langFS, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(path, ".tmpl") {
			return err
		}

		// TODO: error handling
		b, _ := fs.ReadFile(langFS, path)
		_, err = ts.New(path).Parse(string(b))
		return err
	})

	if err != nil {
		panic(fmt.Errorf("error while parsing templates: %w", err))
	}

	return ts
}

var FuncMap = template.FuncMap{
	"KebabToPascalCase": func(s string) string {
		parts := strings.Split(s, "-")

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
	},
	"ScreamSnakeToPascalCase": func(s string) string {
		parts := strings.Split(strings.ToLower(s), "_")

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
	},
	"RemoveLF": func(s string) string {
		s = strings.ReplaceAll(s, "\r", "")
		return strings.ReplaceAll(s, "\n", "")
	},
	"ToLowerCase": func(s string) string {
		return strings.ToLower(s)
	},
}
