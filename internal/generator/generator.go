package generator

import (
	"context"
	"fmt"
	gogen "openassemblybinder/internal/generator/go"
	"openassemblybinder/internal/service"
)

func GenerateGo(
	packageName string,
	clientName string,
	outputPath string,
	createDir bool,
) error {
	globalData := gogen.GlobalTemplateData{
		PackageName: packageName,
		ClientName:  clientName,
	}

	data := gogen.DefaultTemplateData{
		GlobalTemplateData: globalData,
		Header:             map[string]string{
			// TODO : make headers configurable
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
				// TODO : consider logging errors instead of printing to stdout
				fmt.Printf("Error generating service: %v\n", result.Error)
				continue
			}

			svc := result.Service
			bindData := gogen.BindTemplateData{
				GlobalTemplateData: globalData,
				Service:            svc,
			}

			err = gogen.ExecuteBindTemplate(outputPath, bindData)
			if err != nil {
				fmt.Printf("Error executing bind template for service %s: %v\n", svc.StructName, err)
			}
		}
	}
	return nil
}
