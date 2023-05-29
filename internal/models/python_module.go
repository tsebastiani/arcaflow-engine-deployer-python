package models

import (
	"errors"
)

type PythonModule struct {
	fullModuleName string
	ModuleName     *string
	Repo           *string
	ModuleVersion  *string
}

func NewPythonModule(fullModuleName string) PythonModule {
	return PythonModule{fullModuleName: fullModuleName}
}

func (p *PythonModule) PipPackageName() (*string, error) {
	if p.ModuleName == nil {
		return nil, errors.New("PythonModule structure not initialized")
	}
	return &p.fullModuleName, nil
}
