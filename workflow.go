package infocmdbclient

import (
	"encoding/json"
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

func ParseWorkflowParams(jsonParam string) (params WorkflowParams, err error) {
	jsonErr := json.Unmarshal([]byte(jsonParam), &params)
	if jsonErr != nil {
		return params, jsonErr
	}

	return params, nil
}
