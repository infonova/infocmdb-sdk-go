package infocmdb

// This library is used for the communication with the infoCMDB
//
// The API provided by the infoCMDB has different versions and therefore this library is split up v1/v2.
//
// v1 - custom HTTP API
//
// Is the legacy version which is based on configured sql queries provided via a custom http api
//
// v2 - Restful API
//
// This is the first re-engineering or the api to access core models and services to have native access.
// This api properly handles all permission checks and access to native functions.

import (
	"encoding/json"
	v1 "github.com/infonova/infocmdb-sdk-go/infocmdb/v1/infocmdb"
	v2 "github.com/infonova/infocmdb-sdk-go/infocmdb/v2/infocmdb"
	log "github.com/sirupsen/logrus"
	"os"
)

// Client configuration values.
// Usually taken from data/configs/workflows/infocmdb.yml
type Config struct {
	ApiUrl       string `yaml:"apiUrl"`
	ApiUser      string `yaml:"apiUser"`
	ApiPassword  string `yaml:"apiPassword"`
	ApiKey       string
	CmdbBasePath string `yaml:"CmdbBasePath"`
}

// Input parameters to a workflow.
// These are usually passed encoded as a json string as first process argument.
type WorkflowParams struct {
	Apikey              string `json:"apikey"`
	TriggerType         string `json:"triggerType"`
	WorkflowItemId      int    `json:"workflow_item_id"`
	WorkflowInstanceId  int    `json:"workflow_instance_id"`
	CiId                int    `json:"ciid"`
	CiAttributeId       int    `json:"ciAttributeId"`
	CiRelationId        int    `json:"ciRelationId"`
	CiProjectId         int    `json:"ciProjectId"`
	FileImportHistoryId int    `json:"fileImportHistoryId"`
}

// Client combines connectivity methods for version 1 and 2 of the cmdb
type Client struct {
	v1 *v1.Cmdb
	v2 *v2.Cmdb
}

// User defined workflow function that can be passed to `RunWorkflow`.
type Workflow func(params WorkflowParams, cmdb *Client) (err error)

func init() {
	log.SetLevel(log.InfoLevel)
	if os.Getenv("WORKFLOW_DEBUGGING") == "true" {
		log.SetLevel(log.DebugLevel)
	}
}

// NewClient returns a new cmdb client
func NewClient() *Client {
	return new(Client)
}

// LoadConfig from file in yaml format
func (c *Client) LoadConfig(f string) (err error) {
	c.v1, err = v1.New(f)
	if err != nil {
		return
	}

	c.v2 = v2.New()
	if err = c.v2.LoadConfigFile(f); err != nil {
		return
	}

	return
}

// Executes a workflow.
//
// First a infoCMDB client instance is created and the workflow parameters are parsed.
// The workflow parameters are decoded from the first process argument if available.
// Absence of any process argument will not lead to a failure.
//
// If everything is successful the workflow will be executed with the prepared parameters and client.
//
// Any errors that are returned from the workflow function will be logged and lead to a execution failure.
// Additionally the workflow will be marked as failed when something is printed to Stderr during execution.
func RunWorkflow(workflow Workflow) {
	cmdb := NewClient()
	cmdbClientErr := cmdb.LoadConfig("infocmdb.yml")
	if cmdbClientErr != nil {
		log.Fatalf("Failed to Login: %v", cmdbClientErr)
	}

	params, parseErr := parseParams()
	if parseErr != nil {
		log.Fatal(parseErr)
	}

	workflowErr := workflow(params, cmdb)
	if workflowErr != nil {
		log.Fatal(workflowErr)
	}
}

func parseParams() (params WorkflowParams, err error) {
	if len(os.Args) < 2 {
		log.Fatal("Missing json encoded WorkflowParams as first program argument")
	}

	jsonParam := os.Args[1]
	err = json.Unmarshal([]byte(jsonParam), &params)
	if err != nil {
		return
	}

	return
}
