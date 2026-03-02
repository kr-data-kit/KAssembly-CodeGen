package pygen

import (
	"fmt"
	"path"
)

const (
	clientCodeName        = "client.py"
	modelsCodeName        = "models.py"
	statusCodeName        = "status.py"
	exceptionsCodeName    = "exceptions.py"
	initCodeName          = "__init__.py"
	pyprojectCodeName     = "pyproject.toml"
	pythonVersionCodeName = ".python-version"
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

	models := path.Join(dir, modelsCodeName)
	err = ExecuteTemplate(ModelsTemplate, models, data.GlobalTemplateData)
	if err != nil {
		return fmt.Errorf("failed to execute models template: %w", err)
	}

	status := path.Join(dir, statusCodeName)
	err = ExecuteTemplate(StatusTemplate, status, data.GlobalTemplateData)
	if err != nil {
		return fmt.Errorf("failed to execute status template: %w", err)
	}

	exceptions := path.Join(dir, exceptionsCodeName)
	err = ExecuteTemplate(ExceptionsTemplate, exceptions, data.GlobalTemplateData)
	if err != nil {
		return fmt.Errorf("failed to execute exceptions template: %w", err)
	}

	init := path.Join(dir, initCodeName)
	err = ExecuteTemplate(InitTemplate, init, data.GlobalTemplateData)
	if err != nil {
		return fmt.Errorf("failed to execute init template: %w", err)
	}

	pyproject := path.Join(dir, pyprojectCodeName)
	err = ExecuteTemplate(PyprojectTemplate, pyproject, data.GlobalTemplateData)
	if err != nil {
		return fmt.Errorf("failed to execute pyproject template: %w", err)
	}

	pythonVersion := path.Join(dir, pythonVersionCodeName)
	err = ExecuteTemplate(PythonVersionTemplate, pythonVersion, data.GlobalTemplateData)
	if err != nil {
		return fmt.Errorf("failed to execute python-version template: %w", err)
	}

	return nil
}
