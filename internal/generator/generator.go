package generator

import (
	"context"
	"fmt"
	"log/slog"
	gogen "openassemblybinder/internal/generator/go"
	pygen "openassemblybinder/internal/generator/python"
	"openassemblybinder/internal/service"
	"os"
	"path/filepath"
	"sort"
)

const (
	RepositoryURL = "https://github.com/kr-data-kit/KAssembly-CodeGen"
)

func GenerateGo(
	packageName string,
	clientName string,
	outputPath string,
	createDir bool,
) error {
	globalData := gogen.GlobalTemplateData{
		PackageName:   packageName,
		ClientName:    clientName,
		RepositoryURL: RepositoryURL,
	}

	data := gogen.DefaultTemplateData{
		GlobalTemplateData: globalData,
		Header: map[string]string{
			"Content-Type": "application/json",
			"Host":         "open.assembly.go.kr",
			"User-Agent":   "Mozilla/5.0",
		},
	}

	err := gogen.ExecuteDefaultTemplate(outputPath, data)
	if err != nil {
		return fmt.Errorf("failed to execute default template: %v", err)
	}

	// TODO: add ctx config
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	services, err := service.GenerateServices(ctx)
	if err != nil {
		return fmt.Errorf("failed to generate services: %v", err)
	}

	running := true

	for running {
		select {
		case <-ctx.Done():
			running = false
			continue
		case result, ok := <-services:
			if !ok {
				// channel closed, all services processed
				running = false
				continue
			}
			if result.Error != nil {
				slog.Error("Error generating service", "error", result.Error)
				continue
			}

			svc := result.Service
			bindData := gogen.BindTemplateData{
				GlobalTemplateData: globalData,
				Service:            svc,
			}

			err = gogen.ExecuteBindTemplate(outputPath, bindData)
			if err != nil {
				slog.Error("Error executing bind template", "service", svc.StructName, "error", err)
			}
		}
	}
	return nil
}

func GeneratePython(
	packageName string,
	outputPath string,
	createDir bool,
) error {
	// Create endpoints directory
	endpointsDir := filepath.Join(outputPath, "endpoints")
	if createDir {
		err := os.MkdirAll(endpointsDir, 0755)
		if err != nil {
			return fmt.Errorf("failed to create endpoints directory: %v", err)
		}
	}

	globalData := pygen.GlobalTemplateData{
		PackageName:   packageName,
		RepositoryURL: RepositoryURL,
	}

	data := pygen.DefaultTemplateData{
		GlobalTemplateData: globalData,
		Header: map[string]string{
			"Content-Type": "application/json",
			"Host":         "open.assembly.go.kr",
			"User-Agent":   "Mozilla/5.0",
		},
	}

	err := pygen.ExecuteDefaultTemplate(outputPath, data)
	if err != nil {
		return fmt.Errorf("failed to execute default template: %v", err)
	}

	// TODO: add ctx config
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	services, err := service.GenerateServices(ctx)
	if err != nil {
		return fmt.Errorf("failed to generate services: %v", err)
	}

	running := true
	var allServices []*service.Service

	for running {
		select {
		case <-ctx.Done():
			running = false
			continue
		case result, ok := <-services:
			if !ok {
				// channel closed, all services processed
				running = false
				continue
			}
			if result.Error != nil {
				slog.Error("Error generating service", "error", result.Error)
				continue
			}

			svc := result.Service
			allServices = append(allServices, svc)

			endpointData := pygen.EndpointTemplateData{
				GlobalTemplateData: globalData,
				Service:            svc,
			}

			err = pygen.ExecuteEndpointTemplate(outputPath, endpointData)
			if err != nil {
				slog.Error("Error executing endpoint template", "service", svc.ResponseKey, "error", err)
			}
		}
	}

	// Sort services by ResponseKey
	sort.Slice(allServices, func(i, j int) bool {
		return allServices[i].ResponseKey < allServices[j].ResponseKey
	})

	// Generate endpoints/__init__.py with all services
	initData := pygen.EndpointsInitTemplateData{
		GlobalTemplateData: globalData,
		Services:           allServices,
	}

	err = pygen.ExecuteEndpointsInitTemplate(outputPath, initData)
	if err != nil {
		return fmt.Errorf("failed to execute endpoints init template: %v", err)
	}

	return nil
}
