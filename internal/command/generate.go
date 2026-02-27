package command

import (
	"context"
	"fmt"
	"openassemblybinder/internal/generator"
	"openassemblybinder/internal/service"
	"os"
)

func GenerateCommand(
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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	services, err := service.GenerateServices(ctx)
	if err != nil {
		return fmt.Errorf("failed to generate services: %v", err)
	}

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("code generation cancelled: %v", ctx.Err())
		case result, ok := <-services:
			if !ok {
				// channel closed, all services processed
				goto done
			}
			if result.Error != nil {
				fmt.Printf("Error generating service: %v\n", result.Error)
				continue
			}
			svc := result.Service
			bindData := generator.BindTemplateData{
				GlobalTemplateData: globalData,
				Service:            svc,
			}

			err = generator.ExecuteBindTemplate(outputPath, bindData)
			if err != nil {
				fmt.Printf("Error executing bind template for service %s: %v\n", svc.StructName, err)
			}
		}
	}

done:

	fmt.Println("Code generation completed successfully.")

	return nil
}
