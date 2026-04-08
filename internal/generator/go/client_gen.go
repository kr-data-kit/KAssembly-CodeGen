package gogen

import (
	"context"
	"fmt"
	"kassemblycodegen/internal/endpoint"
	"kassemblycodegen/internal/generator"
	"kassemblycodegen/internal/setting"
	"os"
	"sort"
	"strings"
)

type ClientGenerator struct {
	TemplateEngine generator.TemplateEngine
	Endpoints      []*endpoint.Endpoint
	Setting        setting.GoSetting
}

func NewClientGenerator(set setting.GoSetting) *ClientGenerator {
	return &ClientGenerator{
		Setting:        set,
		TemplateEngine: *generator.NewTemplateEngine(ClientTemplates, set.Path),
		Endpoints:      make([]*endpoint.Endpoint, 0),
	}
}

func (c *ClientGenerator) SetGlobalConfig() error {
	if c.Setting.Path == "" {
		return fmt.Errorf("output path is required")
	}
	if c.Setting.CreateDir {
		if err := os.MkdirAll(c.Setting.Path, 0o755); err != nil {
			return fmt.Errorf("failed to create output directory: %w", err)
		}
	}
	return nil
}

func (c *ClientGenerator) GenStatic() error {
	files := map[string]string{
		"client.go.tmpl":     "client.go",
		"interfaces.go.tmpl": "interfaces.go",
		"models.go.tmpl":     "models.go",
		"paging.go.tmpl":     "paging.go",
		"requester.go.tmpl":  "requester.go",
		"status.go.tmpl":     "status.go",
	}
	if c.Setting.IsMod {
		files["go.mod.tmpl"] = "go.mod"
	}

	data := struct {
		RepositoryURL string
		PackageName   string
		Header        map[string]string
	}{
		RepositoryURL: generator.RepositoryURL,
		PackageName:   c.Setting.PackageName,
		Header: map[string]string{
			"Content-Type": "application/json",
			"Host":         "open.assembly.go.kr",
			"User-Agent":   "Mozilla/5.0",
		},
	}

	for name, path := range files {
		err := c.TemplateEngine.Execute(name, path, data)
		if err != nil {
			return fmt.Errorf("failed to execute template %s: %w", name, err)
		}
	}

	return nil
}

func (c *ClientGenerator) GenEndpoint(ctx context.Context, svc *endpoint.Endpoint) error {
	if svc == nil {
		return nil
	}

	data := struct {
		RepositoryURL string
		PackageName   string
		Service       *endpoint.Endpoint
	}{
		RepositoryURL: generator.RepositoryURL,
		PackageName:   c.Setting.PackageName,
		Service:       svc,
	}

	fileName := fmt.Sprintf("bind-%s.go", strings.ToLower(svc.ResponseKey))
	if err := c.TemplateEngine.Execute("endpoint.tmpl", fileName, data); err != nil {
		return fmt.Errorf("failed to execute endpoint template for %s: %w", svc.ResponseKey, err)
	}

	c.Endpoints = append(c.Endpoints, svc)
	return nil
}

func (c *ClientGenerator) GenFinal() error {
	if len(c.Endpoints) == 0 {
		return nil
	}
	sort.Slice(c.Endpoints, func(i, j int) bool {
		return c.Endpoints[i].ResponseKey < c.Endpoints[j].ResponseKey
	})
	return nil
}
