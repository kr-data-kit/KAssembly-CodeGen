package pygen

import (
	"fmt"
	"kassemblycodegen/internal/service"
	"path"
	"strings"
)

type EndpointTemplateData struct {
	GlobalTemplateData
	Service *service.Service
}

type EndpointsInitTemplateData struct {
	GlobalTemplateData
	Services []*service.Service
}

func ExecuteEndpointTemplate(dir string, data EndpointTemplateData) error {
	// Create endpoints subdirectory path
	endpointsDir := path.Join(dir, "endpoints")

	// Generate endpoint module file
	endpoint := path.Join(endpointsDir, createFileName(data.Service.ResponseKey))
	err := ExecuteTemplate(EndpointTemplate, endpoint, data)
	if err != nil {
		return fmt.Errorf("failed to execute endpoint template: %w", err)
	}
	return nil
}

func ExecuteEndpointsInitTemplate(dir string, data EndpointsInitTemplateData) error {
	// Create endpoints subdirectory path
	endpointsDir := path.Join(dir, "endpoints")

	// Generate endpoints __init__.py
	endpointsInit := path.Join(endpointsDir, "__init__.py")
	err := ExecuteTemplate(EndpointInitTemplate, endpointsInit, data)
	if err != nil {
		return fmt.Errorf("failed to execute endpoints init template: %w", err)
	}
	return nil
}

func createFileName(responseKey string) string {
	// Convert ALLBILL to allbill.py
	return strings.ToLower(responseKey) + ".py"
}
