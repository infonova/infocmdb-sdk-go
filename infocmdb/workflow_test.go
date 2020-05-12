package infocmdb

import (
	"os"
	"reflect"
	"testing"
)

func Test_parseParams(t *testing.T) {
	type osArgs []string
	tests := []struct {
		name       string
		osArgs     osArgs
		wantParams WorkflowParams
		wantErr    bool
	}{
		{
			name:       "no os args",
			osArgs:     osArgs{},
			wantParams: WorkflowParams{},
			wantErr:    true,
		},
		{
			name: "only one os arg",
			osArgs: osArgs{
				"test_workflow.go",
			},
			wantParams: WorkflowParams{},
			wantErr:    true,
		},
		{
			name: "empty",
			osArgs: osArgs{
				"test_workflow.go",
				"{}",
			},
			wantParams: WorkflowParams{},
			wantErr:    false,
		},
		{
			name: "manual trigger",
			osArgs: osArgs{
				"test_workflow.go",
				`
{
  "apikey": "6a71bbd97df21cba2dc26d25ac68c5",
  "triggerType": "manual",
  "user_id": "1",
  "workflow_item_id": 32518,
  "workflow_instance_id": 32518
}
`,
			},
			wantParams: WorkflowParams{
				Apikey:              "6a71bbd97df21cba2dc26d25ac68c5",
				TriggerType:         "manual",
				WorkflowItemId:      32518,
				WorkflowInstanceId:  32518,
				CiId:                0,
				CiAttributeId:       0,
				CiRelationId:        0,
				CiProjectId:         0,
				FileImportHistoryId: 0,
			},
			wantErr: false,
		},
		{
			name: "time trigger",
			osArgs: osArgs{
				"test_workflow.go",
				`
{
  "apikey": "6a71bbd97df21cba2dc26d25ac68c5",
  "triggerType": "time",
  "user_id": "1",
  "workflow_item_id": 32518,
  "workflow_instance_id": 32518
}
`,
			},
			wantParams: WorkflowParams{
				Apikey:              "6a71bbd97df21cba2dc26d25ac68c5",
				TriggerType:         "time",
				WorkflowItemId:      32518,
				WorkflowInstanceId:  32518,
				CiId:                0,
				CiAttributeId:       0,
				CiRelationId:        0,
				CiProjectId:         0,
				FileImportHistoryId: 0,
			},
			wantErr: false,
		},
		{
			name: "ci_attribute_create trigger",
			osArgs: osArgs{
				"test_workflow.go",
				`
{
  "apikey": "626409876d3de140670e0325850215",
  "ciAttributeId": 167698,
  "triggerType": "ci_attribute_create",
  "ciid": 14103,
  "user_id": 0,
  "workflow_item_id": 32511,
  "workflow_instance_id": 32511
}
`,
			},
			wantParams: WorkflowParams{
				Apikey:              "626409876d3de140670e0325850215",
				TriggerType:         "ci_attribute_create",
				WorkflowItemId:      32511,
				WorkflowInstanceId:  32511,
				CiId:                14103,
				CiAttributeId:       167698,
				CiRelationId:        0,
				CiProjectId:         0,
				FileImportHistoryId: 0,
			},
			wantErr: false,
		},
		{
			name: "ciid string",
			osArgs: osArgs{
				"test_workflow.go",
				`
{
  "apikey": "626409876d3de140670e0325850215",
  "ciAttributeId": 167698,
  "triggerType": "ci_attribute_create",
  "ciid": "14103",
  "user_id": 0,
  "workflow_item_id": 32511,
  "workflow_instance_id": 32511
}
`,
			},
			wantParams: WorkflowParams{
				Apikey:              "626409876d3de140670e0325850215",
				TriggerType:         "ci_attribute_create",
				WorkflowItemId:      32511,
				WorkflowInstanceId:  32511,
				CiId:                14103,
				CiAttributeId:       167698,
				CiRelationId:        0,
				CiProjectId:         0,
				FileImportHistoryId: 0,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Args = tt.osArgs
			gotParams, err := parseParams()
			if (err != nil) != tt.wantErr {
				t.Errorf("parseParams() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotParams, tt.wantParams) {
				t.Errorf("parseParams() gotParams = %v, want %v", gotParams, tt.wantParams)
			}
		})
	}
}
