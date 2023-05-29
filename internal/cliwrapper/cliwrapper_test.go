package cliwrapper

import (
	"fmt"
	"go.arcalot.io/assert"
	"go.arcalot.io/log/v2"
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
	module := "arcaflow-plugin-template-python@git+https://github.com/tsebastiani/arcaflow-plugin-template-python.git"
	workDir := "/tmp"
	pythonPath := "/usr/bin/python3.9"
	logger := log.NewTestLogger(t)
	python := NewCliWrapper(pythonPath, workDir, logger, false)
	err := pullModule(python, module, workDir, t)
	assert.NoError(t, err)
}

func TestImageExists(t *testing.T) {
	module := "arcaflow-plugin-template-python@git+https://github.com/tsebastiani/arcaflow-plugin-template-python.git"
	workDir := "/tmp"
	pythonPath := "/usr/bin/python3.9"
	logger := log.NewTestLogger(t)
	python := NewCliWrapper(pythonPath, workDir, logger, false)
	removeModuleIfExists(module, python, t)
	exists, err := python.ModuleExists(module)
	assert.Nil(t, err)
	assert.Equals(t, *exists, false)
	err = pullModule(python, module, workDir, t)
	assert.NoError(t, err)
	exists, err = python.ModuleExists(module)
	assert.NoError(t, err)
	assert.Equals(t, *exists, true)
}

func TestImageFormatValidation(t *testing.T) {
	moduleGitNoCommit := "arcaflow-plugin-template-python@git+https://github.com/tsebastiani/arcaflow-plugin-template-python.git"
	moduleGitCommit := "arcaflow-plugin-template-python@git+https://github.com/tsebastiani/arcaflow-plugin-template-python.git@8e43b657db73929d6f8ccb893f059bb67658523f"
	moduleWrongFormat := "https://arcalot.io"
	wrongFormatMessage := "wrong module name format, please use <module-name>@git+<repo_url>[@<commit_sha>]"
	workDir := "/tmp"
	pythonPath := "/usr/bin/python3.9"
	logger := log.NewTestLogger(t)
	wrapperGit := NewCliWrapper(pythonPath, workDir, logger, false)

	// happy path
	path, err := wrapperGit.GetModulePath(moduleGitCommit)
	assert.NoError(t, err)
	assert.Equals(
		t,
		*path,
		fmt.Sprintf("%s/arcaflow-plugin-template-python_8e43b657db73929d6f8ccb893f059bb67658523f", workDir),
	)

	path, err = wrapperGit.GetModulePath(moduleGitNoCommit)
	assert.NoError(t, err)
	assert.Equals(t, *path, fmt.Sprintf("%s/arcaflow-plugin-template-python_latest", workDir))

	_, err = wrapperGit.GetModulePath(moduleWrongFormat)
	assert.Error(t, err)
	assert.Equals(t, err.Error(), wrongFormatMessage)

}
