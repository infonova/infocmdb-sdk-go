package infocmdb

import (
	"encoding/json"
	"regexp"
)

const (
	WORKFLOW_TRIGGER_TYPE_CI_CREATE = "ci_create"
	WORKFLOW_TRIGGER_TYPE_CI_UPDATE = "ci_update"
	WORKFLOW_TRIGGER_TYPE_CI_DELETE = "ci_delete"

	WORKFLOW_TRIGGER_TYPE_CI_TYPE_CHANGE_UPDATE = "ci_type_change_update"

	WORKFLOW_TRIGGER_TYPE_CI_ATTRIBUTE_CREATE = "ci_attribute_create"
	WORKFLOW_TRIGGER_TYPE_CI_ATTRIBUTE_UPDATE = "ci_attribute_update"
	WORKFLOW_TRIGGER_TYPE_CI_ATTRIBUTE_DELETE = "ci_attribute_delete"

	WORKFLOW_TRIGGER_TYPE_CI_RELATION_CREATE = "ci_relation_create"
	WORKFLOW_TRIGGER_TYPE_CI_RELATION_DELETE = "ci_relation_delete"

	WORKFLOW_TRIGGER_TYPE_CI_PROJECT_CREATE = "ci_project_create"
	WORKFLOW_TRIGGER_TYPE_CI_PROJECT_DELETE = "ci_project_delete"

	WORKFLOW_TRIGGER_TYPE_FILEIMPORT_BEFORE           = "fileimport_before"
	WORKFLOW_TRIGGER_TYPE_FILEIMPORT_AFTER            = "fileimport_after"
	WORKFLOW_TRIGGER_TYPE_FILEIMPORT_BEFORE_AND_AFTER = "fileimport_before_and_after"
)

type WorkflowParams struct {
	Apikey              string `json:"apikey"`
	TriggerType         string `json:"triggerType"`
	WorkflowItemId      int    `json:"workflow_item_id,string"`
	WorkflowInstanceId  int    `json:"workflow_instance_id,string"`
	CiId                int    `json:"ciid,string"`
	CiAttributeId       int    `json:"ciAttributeId,string"`
	CiRelationId        int    `json:"ciRelationId,string"`
	CiProjectId         int    `json:"ciProjectId,string"`
	FileImportHistoryId int    `json:"fileImportHistoryId,string"`
}

var reJsonParam = regexp.MustCompile(`^\s*'(.*?)\s*'$`)

func ParseWorkflowParams(jsonParam string) (params WorkflowParams, err error) {
	if reJsonParam.MatchString(jsonParam) {
		jsonParam = reJsonParam.FindStringSubmatch(jsonParam)[1]
	}
	jsonErr := json.Unmarshal([]byte(jsonParam), &params)
	if jsonErr != nil {
		return params, jsonErr
	}

	return params, nil
}
