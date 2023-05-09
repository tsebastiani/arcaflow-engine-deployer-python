package models

import (
	"errors"
	"fmt"
	python "go.flow.arcalot.io/pythondeployer/internal/config"
)

type PythonModule struct {
	fullModuleName string
	ModuleName     *string
	Repo           *string
	ModuleVersion  *string
	moduleSource   python.ModuleSource
}

func NewPythonModule(source python.ModuleSource, fullModuleName string) PythonModule {
	return PythonModule{moduleSource: source, fullModuleName: fullModuleName}
}

func (p *PythonModule) PipPackageName() (*string, error) {
	if p.ModuleName == nil {
		return nil, errors.New("PythonModule structure not initialized")
	}
	if p.moduleSource == python.ModuleSourceGit {
		return &p.fullModuleName, nil
	}

	if p.ModuleVersion != nil {
		packageName := fmt.Sprintf("%s==%s", *p.ModuleName, *p.ModuleVersion)
		return &packageName, nil
	}

	return p.ModuleName, nil

}
