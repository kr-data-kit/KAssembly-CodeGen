package gogen

import (
	"fmt"
	"path"
)

const (
	clientCodeName = "client.go"
	commonCodeName = "common.go"
	statusCodeName = "status.go"
	goModCodeName  = "go.mod"
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
	common := path.Join(dir, commonCodeName)
	err = ExecuteTemplate(CommonTemplate, common, data.GlobalTemplateData)
	if err != nil {
		return fmt.Errorf("failed to execute common template: %w", err)
	}
	status := path.Join(dir, statusCodeName)
	err = ExecuteTemplate(StatusTemplate, status, data.GlobalTemplateData)
	if err != nil {
		return fmt.Errorf("failed to execute status template: %w", err)
	}
	// TODO : consider adding go.mod template execution
	return nil
}
