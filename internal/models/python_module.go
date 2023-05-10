package models

import (
	"errors"
	"fmt"
	"github.com/tsebastiani/arcaflow-engine-deployer-python/internal/config"
)

type PythonModule struct {
	fullModuleName string
	ModuleName     *string
	Repo           *string
	ModuleVersion  *string
	moduleSource   config.ModuleSource
}

func NewPythonModule(source config.ModuleSource, fullModuleName string) PythonModule {
	return PythonModule{moduleSource: source, fullModuleName: fullModuleName}
}

func (p *PythonModule) PipPackageName() (*string, error) {
	if p.ModuleName == nil {
		return nil, errors.New("PythonModule structure not initialized")
	}
	if p.moduleSource == config.ModuleSourceGit {
		return &p.fullModuleName, nil
	}

	if p.ModuleVersion != nil {
		packageName := fmt.Sprintf("%s==%s", *p.ModuleName, *p.ModuleVersion)
		return &packageName, nil
	}

	return p.ModuleName, nil

}
