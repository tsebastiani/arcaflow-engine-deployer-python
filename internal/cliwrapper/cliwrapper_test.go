package cliwrapper

import (
	"fmt"
	"go.arcalot.io/assert"
	"go.arcalot.io/log"
	"go.flow.arcalot.io/pythondeployer/internal/config"
	"os"
	"testing"
)

func removeModuleIfExists(module string, python CliWrapper, t *testing.T) {
	modulePath, err := python.GetModulePath(module)
	assert.Nil(t, err)
	if _, err := os.Stat(*modulePath); !os.IsNotExist(err) {
		os.RemoveAll(*modulePath)
	}
}

func pullModule(python CliWrapper, module string, workDir string, t *testing.T) error {
	removeModuleIfExists(module, python, t)
	err := python.PullModule(module)
	if err != nil {
		return err
	}
	return nil
}

func TestPullImage(t *testing.T) {
	module := "krkn-lib-kubernetes@0.1.0"
	workDir := "/tmp"
	pythonPath := "/usr/bin/python3.9"
	logger := log.NewTestLogger(t)
	python := NewCliWrapper(pythonPath, workDir, config.ModuleSourcePypi, logger)
	pullModule(python, module, workDir, t)
}

func TestImageExists(t *testing.T) {
	module := "krkn-lib-kubernetes@0.1.0"
	workDir := "/tmp"
	pythonPath := "/usr/bin/python3.9"
	logger := log.NewTestLogger(t)
	python := NewCliWrapper(pythonPath, workDir, config.ModuleSourcePypi, logger)
	removeModuleIfExists(module, python, t)
	exists, err := python.ModuleExists(module)
	assert.Nil(t, err)
	assert.Equals(t, *exists, false)
	pullModule(python, module, workDir, t)
	exists, err = python.ModuleExists(module)
	assert.Nil(t, err)
	assert.Equals(t, *exists, true)
}

func TestImageFormatValidation(t *testing.T) {
	modulePypiVersion := "krkn-lib-kubernetes@0.1.0"
	modulePypiNoVersion := "krkn-lib-kubernetes"
	moduleGitNoCommit := "arcaflow-plugin-template-python@git+https://github.com/tsebastiani/arcaflow-plugin-template-python.git"
	moduleGitCommit := "arcaflow-plugin-template-python@git+https://github.com/tsebastiani/arcaflow-plugin-template-python.git@8e43b657db73929d6f8ccb893f059bb67658523f"
	moduleWrongFormat := "https://arcalot.io"

	wrongFormatMessage := "wrong module name format"

	gitMismatchMessage := "you're using a pip module name " +
		"format using the deployer in git mode, " +
		"please change the deployer configuration"

	pypiMismatchMessage := "you're using a git module name " +
		"format using the deployer in pipy mode, " +
		"please change the deployer configuration"

	workDir := "/tmp"
	pythonPath := "/usr/bin/python3.9"
	logger := log.NewTestLogger(t)
	wrapperPypi := NewCliWrapper(pythonPath, workDir, config.ModuleSourcePypi, logger)
	wrapperGit := NewCliWrapper(pythonPath, workDir, config.ModuleSourceGit, logger)

	// happy path
	path, err := wrapperPypi.GetModulePath(modulePypiVersion)
	assert.NoError(t, err)
	assert.Equals(t, *path, fmt.Sprintf("%s/krkn-lib-kubernetes_0.1.0", workDir))
	path, err = wrapperPypi.GetModulePath(modulePypiNoVersion)
	assert.NoError(t, err)
	assert.Equals(t, *path, fmt.Sprintf("%s/krkn-lib-kubernetes_latest", workDir))

	path, err = wrapperGit.GetModulePath(moduleGitCommit)
	assert.NoError(t, err)
	assert.Equals(
		t,
		*path,
		fmt.Sprintf("%s/arcaflow-plugin-template-python_8e43b657db73929d6f8ccb893f059bb67658523f", workDir),
	)

	path, err = wrapperGit.GetModulePath(moduleGitNoCommit)
	assert.NoError(t, err)
	assert.Equals(t, *path, fmt.Sprintf("%s/arcaflow-plugin-template-python_latest", workDir))

	path, err = wrapperGit.GetModulePath(moduleWrongFormat)
	assert.Error(t, err)
	path, err = wrapperPypi.GetModulePath(moduleWrongFormat)
	assert.Error(t, err)
	assert.Equals(t, err.Error(), wrongFormatMessage)

	// Pypi source with git module name
	path, err = wrapperPypi.GetModulePath(moduleGitCommit)
	assert.Error(t, err)
	assert.Equals(t, err.Error(), pypiMismatchMessage)

	path, err = wrapperPypi.GetModulePath(moduleGitNoCommit)
	assert.Error(t, err)
	assert.Equals(t, err.Error(), pypiMismatchMessage)

	// Git source with pypi module name

	path, err = wrapperGit.GetModulePath(modulePypiNoVersion)
	assert.Error(t, err)
	assert.Equals(t, err.Error(), gitMismatchMessage)

	path, err = wrapperGit.GetModulePath(modulePypiVersion)
	assert.Error(t, err)
	assert.Equals(t, err.Error(), gitMismatchMessage)

}
