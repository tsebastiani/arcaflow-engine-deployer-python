# Arcaflow engine deployer python

This library is an implementation of the [arcaflow deployer interface](https://github.com/arcalot/arcaflow-engine-deployer) that uses the podman CLI.

## Usage

### Deployer configuration (config.yaml)
(full configuration example)
```
deployer:
  type: python
  # Optional Fields
  pythonPath: /usr/bin/python3.9
  workdir: /tmp
  modulePullPolicy: Always | IfNotPresent
  moduleSource: Git | Pypi
```
- `pythonPath` (_optional_, default `/usr/bin/python`)
  - Path to the python interpreter binary 
- `workdir` (_optional_, default `/tmp`)
  - folder where the virtual environments of every single plugin are stored. 
    Setting `modulePullPolicy` as `IfNotPresent` the workdir will work as a cache 
    and will speed up the workflow runs.
- `modulePullPolicy` (_optional_, default `IfNotPresent`)
  - `IfNotPresent`: will check in the `workdir` path if the requested module
    At the requested version has been already pulled
  - `Always`: will always pull the module, if already present, will delete the previous
    version and will pull it again.
- `moduleSource` (_optional_, default `Pypi`)
  - `Pypi`: the module will be pulled from pypi.org.
    The plugin name must be in the format `<pypi_module_name>[@<version>]`,
    if version is omitted the latest version will be pulled
  - `Git` : the module will be pulled from a Git repository containing a plugin.
    The plugin name must be in the format `<module_name>@git+<repo_url>[@git_commit_sha]`,
    if git commit sha is omitted, the HEAD of main will be pulled

## Worfklows (workflow.yaml)
The main difference in the workflow syntax is that instead of passing a container image
as plugin (like the podman, docker and kubernetes deployer) must be passed a python module
either in the `Git` or in the `Pypi` format as previously mentioned.

Module Name Format:
- `Pypi`: `<pypi_module_name>[@<version>]`
- `Git`: `<module_name>@git+<repo_url>[@git_commit_sha]`

Example `Git` source workflow
```
steps:
  kill_pod:
    plugin: arcaflow-plugin-kill-pod@git+https://github.com/redhat-chaos/arcaflow-plugin-kill-pod.git@a34551a4aa68d822ba54f338148ca6e6a28c493b
    step: kill-pods
    input:
    ...
```

Example `PyPi` source workflow
```
steps:
  kill_pod:
    plugin: arcaflow-plugin-kill-pod@0.1.0
    step: kill-pods
    input:
    ....
```

