package gogen

import (
	"fmt"
	"path"
)

const (
	clientCodeName     = "client.go"
	pagingCodeName     = "paging.go"
	modelsCodeName     = "models.go"
	interfacesCodeName = "interfaces.go"
	requesterCodeName  = "requester.go"
	statusCodeName     = "status.go"
	goModCodeName      = "go.mod"
)

type DefaultTemplateData struct {
	GlobalTemplateData
	Header map[string]string // for default headers (client file)
}

func ExecuteDefaultTemplate(dir string, data DefaultTemplateData) error {
	client := path.Join(dir, clientCodeName)
	err := ExecuteTemplate(ClientTemplate, client, data)
	if err != nil {
		return fmt.Errorf("failed to execute client template: %w", err)
	}
	paging := path.Join(dir, pagingCodeName)
	err = ExecuteTemplate(PagingTemplate, paging, data.GlobalTemplateData)
	if err != nil {
		return fmt.Errorf("failed to execute paging template: %w", err)
	}
	models := path.Join(dir, modelsCodeName)
	err = ExecuteTemplate(ModelsTemplate, models, data.GlobalTemplateData)
	if err != nil {
		return fmt.Errorf("failed to execute models template: %w", err)
	}
	interfaces := path.Join(dir, interfacesCodeName)
	err = ExecuteTemplate(InterfacesTemplate, interfaces, data.GlobalTemplateData)
	if err != nil {
		return fmt.Errorf("failed to execute interfaces template: %w", err)
	}
	requester := path.Join(dir, requesterCodeName)
	err = ExecuteTemplate(RequesterTemplate, requester, data.GlobalTemplateData)
	if err != nil {
		return fmt.Errorf("failed to execute requester template: %w", err)
	}
	status := path.Join(dir, statusCodeName)
	err = ExecuteTemplate(StatusTemplate, status, data.GlobalTemplateData)
	if err != nil {
		return fmt.Errorf("failed to execute status template: %w", err)
	}
	// TODO : consider adding go.mod template execution
	return nil
}
