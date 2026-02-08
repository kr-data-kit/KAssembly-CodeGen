package command

import (
	"context"
	"fmt"
	"openassemblybind/internal/generator"
	"openassemblybind/internal/service"
	"os"
)

func GenerateCommand(
	key string,
	packageName string,
	clientName string,
	outputPath string,
	createDir bool,
) error {
	// just for test, print the parameters
	fmt.Println("Generating code with the following parameters:")
	fmt.Println("Package Name:", packageName)
	fmt.Println("Client Name :", clientName)
	fmt.Println("Output Path :", outputPath)
	fmt.Println("Create Dir  :", createDir)
	fmt.Println("")

	// check output directory existence (not creating directory here)
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		if createDir {
			err := os.MkdirAll(outputPath, os.ModePerm)
			if err != nil {
				return fmt.Errorf("failed to create output directory: %v", err)
			}
		} else {
			return fmt.Errorf("output directory does not exist: %s", outputPath)
		}
	}

	globalData := generator.GlobalTemplateData{
		PackageName: packageName,
		ClientName:  clientName,
	}

	data := generator.DefaultTemplateData{
		GlobalTemplateData: globalData,
		Header: map[string]string{
			// TODO : make headers configurable
			"Content-Type": "application/json",
			"Host":         "open.assembly.go.kr",
			"User-Agent":   "Mozilla/5.0",
		},
	}

	err := generator.ExecuteDefaultTemplate(outputPath, data)
	if err != nil {
		return fmt.Errorf("failed to execute default template: %v", err)
	}

	// TODO: add ctx config
	ctx := context.Background()
	summaries, err := service.FetchServiceSummaries(
		ctx,
		key,
	)
	if err != nil {
		return fmt.Errorf("failed to fetch service summaries: %v", err)
	}

	for _, summary := range summaries {
		svc, err := service.CreateService(
			ctx,
			summary,
			service.CreateServiceOptions{
				// TODO : make this configurable
				CreateServiceEn: true,
			},
		)
		if err != nil {
			return fmt.Errorf("failed to create service for %s(%s): %v", summary.ID, summary.Title, err)
		}
		bindData := generator.BindTemplateData{
			GlobalTemplateData: globalData,
			Service:            svc,
		}

		err = generator.ExecuteBindTemplate(outputPath, bindData)
		if err != nil {
			return fmt.Errorf("failed to execute bind template for %s: %v", svc.StructName, err)
		}
	}

	fmt.Println("Code generation completed successfully.")

	return nil
}
