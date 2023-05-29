package config

type Config struct {
	PythonPath                  string           `json:"pythonPath"`
	WorkDir                     string           `json:"workdir"`
	ModulePullPolicy            ModulePullPolicy `json:"modulePullPolicy"`
	ModuleSource                ModuleSource     `json:"moduleSource"`
	OverrideModuleCompatibility bool             `json:"overrideModuleCompatibility"`
}

type ModulePullPolicy string
type ModuleSource string

const (
	// ModulePullPolicyAlways means that the module will be pulled for every workflow run.
	ModulePullPolicyAlways ModulePullPolicy = "Always"
	// ModulePullPolicyIfNotPresent means the image will be pulled if the module is not present locally
	ModulePullPolicyIfNotPresent ModulePullPolicy = "IfNotPresent"
	ModuleSourcePypi             ModuleSource     = "Pypi"
	ModuleSourceGit              ModuleSource     = "Git"
)
