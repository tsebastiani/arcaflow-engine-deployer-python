package pythondeployer

import (
	"go.arcalot.io/log"
	"go.flow.arcalot.io/deployer"
	"go.flow.arcalot.io/pluginsdk/schema"
	"go.flow.arcalot.io/pythondeployer/internal/cliwrapper"
	"go.flow.arcalot.io/pythondeployer/internal/config"
)

// NewFactory creates a new factory for the Docker deployer.
func NewFactory() deployer.ConnectorFactory[*config.Config] {
	return &factory{}
}

type factory struct {
}

func (f factory) ID() string {
	return "docker"
}

func (f factory) ConfigurationSchema() *schema.TypedScopeSchema[*config.Config] {
	return Schema
}

func (f factory) Create(config *config.Config, logger log.Logger) (deployer.Connector, error) {
	python := cliwrapper.NewCliWrapper(config.PythonPath, config.WorkDir, config.ModuleSource, logger)
	return &Connector{
		config:           config,
		logger:           logger,
		pythonCliWrapper: python,
	}, nil
}