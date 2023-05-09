package cliwrapper

import (
	"io"
)

type CliWrapper interface {
	ModuleExists(fullModuleName string) (*bool, error)
	PullModule(fullModuleName string) error
	Deploy(image string) (io.WriteCloser, io.ReadCloser, error)
	KillAndClean(moduleName string) error
	GetModulePath(fullModuleName string) (*string, error)
}
