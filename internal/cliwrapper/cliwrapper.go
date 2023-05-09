package cliwrapper

import (
	"bytes"
	"errors"
	"fmt"
	python "go.flow.arcalot.io/pythondeployer/internal/config"
	"go.flow.arcalot.io/pythondeployer/internal/models"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"go.arcalot.io/log"
)

type cliWrapper struct {
	pythonFullPath string
	workDir        string
	moduleSource   python.ModuleSource
	deployCommand  *exec.Cmd
	logger         log.Logger
}

func NewCliWrapper(pythonFullPath string, workDir string, moduleSource python.ModuleSource, logger log.Logger) CliWrapper {
	return &cliWrapper{
		pythonFullPath: pythonFullPath,
		logger:         logger,
		workDir:        workDir,
		moduleSource:   moduleSource,
	}
}

func parseModuleNameGit(fullName string, module *models.PythonModule) {
	nameSourceVersion := strings.Split(fullName, "@")
	source := strings.Replace(nameSourceVersion[1], "git+", "", 1)
	(*module).ModuleName = &nameSourceVersion[0]
	(*module).Repo = &source
	if len(nameSourceVersion) == 3 {
		(*module).ModuleVersion = &nameSourceVersion[2]
	}
}

func parseModuleNamePip(fullName string, module *models.PythonModule) {
	nameAndVersion := strings.Split(fullName, "@")
	(*module).ModuleName = &nameAndVersion[0]
	if len(nameAndVersion) == 2 {
		(*module).ModuleVersion = &nameAndVersion[1]
	}
}

func parseModuleName(fullName string, moduleSource python.ModuleSource) (*models.PythonModule, error) {
	pythonModule := models.NewPythonModule(moduleSource, fullName)
	pypiRegex := "^([a-zA-Z0-9]+[_,-]*)+$|^([a-zA-Z0-9]+[_,-]*)+@[a-zA-Z0-9\\.]+$"
	gitRegex := "^([a-zA-Z0-9]+[-_]*)+@git\\+http[s]{0,1}:\\/\\/([a-zA-Z0-9]+[-.\\/]*)+(@[a-z0-9]+)*$"
	matchPypi, _ := regexp.MatchString(pypiRegex, fullName)
	matchGit, _ := regexp.MatchString(gitRegex, fullName)

	if matchPypi && moduleSource == python.ModuleSourceGit {
		return nil, errors.New("you're using a pip module name " +
			"format using the deployer in git mode, " +
			"please change the deployer configuration")
	}
	if matchGit && moduleSource == python.ModuleSourcePypi {
		return nil, errors.New("you're using a git module name " +
			"format using the deployer in pip mode, " +
			"please change the deployer configuration")
	}
	if !matchGit && !matchPypi {
		return nil, errors.New("wrong module name format")
	}

	if matchGit {
		parseModuleNameGit(fullName, &pythonModule)
	} else {
		parseModuleNamePip(fullName, &pythonModule)
	}

	return &pythonModule, nil
}

func (p *cliWrapper) GetModulePath(fullModuleName string) (*string, error) {
	pythonModule, err := parseModuleName(fullModuleName, p.moduleSource)
	if err != nil {
		return nil, err
	}
	modulePath := ""
	if pythonModule.ModuleVersion != nil {
		modulePath = fmt.Sprintf("%s/%s_%s", p.workDir, *pythonModule.ModuleName, *pythonModule.ModuleVersion)
	} else {
		modulePath = fmt.Sprintf("%s/%s_latest", p.workDir, *pythonModule.ModuleName)
	}
	return &modulePath, err
}

func (p *cliWrapper) ModuleExists(fullModuleName string) (*bool, error) {
	moduleExists := false
	modulePath, err := p.GetModulePath(fullModuleName)
	if err != nil {
		return nil, err
	}

	if _, err := os.Stat(*modulePath); os.IsNotExist(err) {
		return &moduleExists, nil
	}

	moduleExists = true
	return &moduleExists, nil
}

func (p *cliWrapper) PullModule(fullModuleName string) error {
	pythonModule, err := parseModuleName(fullModuleName, p.moduleSource)
	if err != nil {
		return err
	}
	module, err := pythonModule.PipPackageName()
	if err != nil {
		return err
	}
	// create module path
	modulePath, err := p.GetModulePath(fullModuleName)
	if err != nil {
		return err
	}
	if err := os.Mkdir(*modulePath, os.ModePerm); err != nil {
		return err
	}

	// create venv

	cmdCreateVenv := exec.Command(p.pythonFullPath, "-m", "venv", "venv")
	cmdCreateVenv.Dir = *modulePath
	var cmdCreateOut bytes.Buffer
	cmdCreateVenv.Stderr = &cmdCreateOut
	if err := cmdCreateVenv.Run(); err != nil {
		return errors.New(cmdCreateOut.String())
	}

	// pull module
	pipPath := fmt.Sprintf("%s/venv/bin/pip", *modulePath)
	cmdPip := exec.Command(pipPath, "install", *module)
	var cmdPipOut bytes.Buffer
	cmdPip.Stderr = &cmdPipOut
	if err := cmdPip.Run(); err != nil {
		return errors.New(cmdPipOut.String())
	}
	return nil
}

func (p *cliWrapper) Deploy(image string) (io.WriteCloser, io.ReadCloser, error) {
	pythonModule, err := parseModuleName(image, p.moduleSource)
	if err != nil {
		return nil, nil, err
	}
	args := []string{"-m"}
	moduleInvokableName := strings.ReplaceAll(*pythonModule.ModuleName, "-", "_")
	args = append(args, moduleInvokableName)
	args = append(args, "--atp")
	venvPath, err := p.GetModulePath(image)
	if err != nil {
		return nil, nil, err
	}
	venvPython := fmt.Sprintf("%s/venv/bin/python", *venvPath)
	p.deployCommand = exec.Command(venvPython, args...) //nolint:gosec
	var stdErrBuff bytes.Buffer
	p.deployCommand.Stderr = &stdErrBuff
	stdin, err := p.deployCommand.StdinPipe()
	if err != nil {
		return nil, nil, err
	}
	stdout, err := p.deployCommand.StdoutPipe()
	if err != nil {
		return nil, nil, err
	}

	if err := p.deployCommand.Start(); err != nil {
		return nil, nil, errors.New(stdErrBuff.String())
	}
	return stdin, stdout, nil
}

func (p *cliWrapper) KillAndClean(containerName string) error {
	p.logger.Infof("killing python process with pid %d", p.deployCommand.Process.Pid)
	err := p.deployCommand.Process.Kill()
	if err != nil {
		return err
	}
	return nil
}
