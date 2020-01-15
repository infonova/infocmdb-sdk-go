# infocmdb-sdk-go

[![Build Status](https://travis-ci.com/infonova/infocmdb-sdk-go.svg?branch=master)](https://travis-ci.com/infonova/infocmdb-sdk-go)

Library for the [infoCMDB](https://github.com/infonova/infocmdb) REST API.\
It provides easy integration with workflows and access to most webserivces and v2 APIs.

## Contents

* [Usage in workflows](#usage-in-workflows)
    * [Workflow script](#workflow-script)
    * [Workflow test](#workflow-test)
* [Recommendation for workflow code](#recommendation-for-workflow-code)
* [Logging](#logging)
* [License](#license)

## Usage in workflows

### Workflow script

The most common use case of this library is made especially easy.

```go
package main

import (
    "github.com/infonova/infocmdb-sdk-go/infocmdb"
    log "github.com/sirupsen/logrus"
)

func main() {
    w := infocmdb.NewWorkflow()
    // If you need to load a non-default CMDB config you first need to set the config file location.
    // Both absolute paths and relative paths (which depend on the WORKFLOW_CONFIG_PATH) are supported.
    // The default value is "infocmdb.yml".
    w.SetConfig("demoCMDB.yml") 
    w.Run(workflow)
}

func workflow(params infocmdb.WorkflowParams, cmdb *infocmdb.Client) (err error) {
    log.Info("Implement the workflow logic here")
    
    // If you need access to the workflow context for ci trigger workflows:
    ctx, err := cmdb.GetWorkflowContext(params.WorkflowInstanceId)
    if err != nil {
        return err
    }

    log.Infof("Old ci data: %+v", ctx.Data.Old)
    log.Infof("New ci data: %+v", ctx.Data.New)

    return
}
```

### Workflow test

Workflow tests are executed prior to compilation for any change.\
They usually test that all required preconditions (CI Types, attributes, Relation Types, ...) exist.\
This prevents common workflow errors when migrating workflows from different stages (e.g. DEV to PROD).

```go
package main

import (
    "github.com/infonova/infocmdb-sdk-go/infocmdb"
    "testing"
)

func TestPreconditions(t *testing.T) {
    w := infocmdb.NewWorkflow()
    w.TestPreconditions(t, infocmdb.Preconditions{
        {infocmdb.TYPE_CI_TYPE, "required_ci_type"},
        {infocmdb.TYPE_ATTRIBUTE, "required_attribute"},
        {infocmdb.TYPE_RELATION, "required_relation"},
    })
}
```

## Recommendation for workflow code

Although all workflow logic could implemented directly in infoCMDB, it is **not** recommended to do so.\
Workflows should be kept simple, otherwise they become difficult to maintain.

If complex workflow logic is required, it should be extracted to a separate library.\
This provides several advantages:
* A separate workflow library in its own repo grants all the VCS benefits like commit history, etc.
* Easy modularization, split packages
* Dependency management via go modules
* Unit tests for separate files and packages, custom integration tests
* Automated test execution via Continuous Integration servers
* Easier collobaration with Pull Requests
* You name it...

## Logging

This library makes heavy use of [logrus](https://github.com/sirupsen/logrus) 
and workflow code (or libraries called from the workflow) are advised to do the same.

The global logger is preconfigured on initialization:
* Logs with level info, debug and trace are written to stdout, all others are written to stderr\
  If anything is written to stderr, the workflow will continue running but will get status `FAILED`
* The [log format](https://github.com/t-tomalak/logrus-easy-formatter) is changed to `[%lvl%] %msg%\n`\
  The timestamp is omitted because it is provided by infoCMDB in a separate column already. 
* Default log level is `INFO` (default) or `DEBUG` depending on the `WORKFLOW_DEBUGGING` environment variable

If you need to see infoCMDB responses, enable debug logging by setting the env variable `WORKFLOW_DEBUGGING` to `true`.

## License

This project is licensed under the Apache License 2.0 - see LICENSE file for details.
