package pythondeployer

import (
	"context"
	"encoding/json"
	"github.com/tsebastiani/arcaflow-engine-deployer-python/internal/cliwrapper"
	"github.com/tsebastiani/arcaflow-engine-deployer-python/internal/config"
	"go.arcalot.io/assert"
	"go.arcalot.io/log/v2"
	"go.flow.arcalot.io/deployer"
	"go.flow.arcalot.io/pluginsdk/atp"
	"os"
	"testing"
)

func getCliWrapper(t *testing.T, overrideCompatibilityCheck bool) cliwrapper.CliWrapper {
	workDir := "/tmp"
	pythonPath := "/usr/bin/python3.9"
	logger := log.NewTestLogger(t)
	return cliwrapper.NewCliWrapper(pythonPath, workDir, logger, overrideCompatibilityCheck)
}

func getConnector(t *testing.T, configJSON string) (deployer.Connector, *config.Config) {
	var config any
	if err := json.Unmarshal([]byte(configJSON), &config); err != nil {
		t.Fatal(err)
	}
	factory := NewFactory()
	schema := factory.ConfigurationSchema()
	unserializedConfig, err := schema.UnserializeType(config)
	assert.NoError(t, err)
	connector, err := factory.Create(unserializedConfig, log.NewTestLogger(t))
	assert.NoError(t, err)
	return connector, unserializedConfig
}

var inOutConfigGitPullAlways = `
{
	"pythonPath":"/usr/bin/python3.9",
	"workdir":"/tmp",
	"modulePullPolicy":"Always"
}
`

var inOutConfigGitPullIfNotPresent = `
{
	"pythonPath":"/usr/bin/python3.9",
	"workdir":"/tmp",
	"modulePullPolicy":"IfNotPresent"
}
`

func TestRunStepGit(t *testing.T) {
	moduleName := "arcaflow-plugin-template-python@git+https://github.com/tsebastiani/arcaflow-plugin-template-python.git@faeffde803696d85756d05afd74dd5bd8c9519e5"
	connector, _ := getConnector(t, inOutConfigGitPullAlways)
	RunStep(t, connector, moduleName)
}

func TestPullPolicies(t *testing.T) {
	moduleName := "arcaflow-plugin-template-python@git+https://github.com/tsebastiani/arcaflow-plugin-template-python.git"
	connectorAlways, _ := getConnector(t, inOutConfigGitPullAlways)
	connectorIfNotPresent, _ := getConnector(t, inOutConfigGitPullIfNotPresent)
	// pull mode Always, venv will be removed if present and pulled again
	RunStep(t, connectorAlways, moduleName)
	// pull mode IfNotPresent, venv will be kept
	RunStep(t, connectorIfNotPresent, moduleName)
	wrapper := getCliWrapper(t, false)
	path, err := wrapper.GetModulePath(moduleName)
	assert.NoError(t, err)
	file, err := os.Stat(*path)
	assert.NoError(t, err)
	// venv path modification time is checked
	startTime := file.ModTime()
	// pull mode Always, venv will be removed if present and pulled again
	RunStep(t, connectorAlways, moduleName)
	file, err = os.Stat(*path)
	assert.NoError(t, err)
	// venv path modification time is checked
	newTime := file.ModTime()
	// new time check must be greater than the first one checked
	assert.Equals(t, newTime.After(startTime), true)
}

func RunStep(t *testing.T, connector deployer.Connector, moduleName string) {
	stepID := "hello-world"
	input := map[string]any{"name": "Tester McTesty"}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	plugin, err := connector.Deploy(context.Background(), moduleName)
	assert.NoError(t, err)
	t.Cleanup(func() {
		assert.NoError(t, plugin.Close())
	})

	atpClient := atp.NewClient(plugin)
	pluginSchema, err := atpClient.ReadSchema()
	assert.NoError(t, err)
	steps := pluginSchema.Steps()
	step, ok := steps[stepID]
	if !ok {
		t.Fatalf("no such step: %s", stepID)
	}

	_, err = step.Input().Unserialize(input)
	assert.NoError(t, err)

	// Executes the step and validates that the output is correct.
	outputID, outputData, _ := atpClient.Execute(ctx, stepID, input)
	assert.Equals(t, outputID, "success")
	assert.NoError(t, err)
	assert.NotNil(t, pluginSchema)
	assert.Equals(t,
		outputData.(map[interface{}]interface{}),
		map[interface{}]interface{}{"message": "Hello, Tester McTesty!"})
}
