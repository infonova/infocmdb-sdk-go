package infocmdb

import (
	"encoding/json"
	"errors"
	log "github.com/sirupsen/logrus"
	"strconv"
)

const (
	WORKFLOW_ENV_APPLICATION_ENV      = "APPLICATION_ENV"
	WORKFLOW_ENV_APPLICATION_PATH     = "APPLICATION_PATH"
	WORKFLOW_ENV_APPLICATION_URL      = "APPLICATION_URL"
	WORKFLOW_ENV_APPLICATION_DATA     = "APPLICATION_DATA"
	WORKFLOW_ENV_APPLICATION_PUBLIC   = "APPLICATION_PUBLIC"
	WORKFLOW_ENV_WORKFLOW_CONFIG_PATH = "WORKFLOW_CONFIG_PATH"
	WORKFLOW_ENV_WORKFLOW_DEBUGGING   = "WORKFLOW_DEBUGGING"

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

type WorkflowContext struct {
	Environment map[string]string `json:"Environment"`
	Ciid        int               `json:"ciid"`
	TriggerType string            `json:"triggerType"`
	Data        Data              `json:"data"`
	UserID      string            `json:"user_id"`
}
type Data struct {
	Old *CiDetail `json:"old"`
	New *CiDetail `json:"new"`
}
type CiDetail struct {
	//Relations  map[int]Relation          `json:"relations"`
	Projects   map[int]Project           `json:"projects"`
	CiTypeID   string                    `json:"ciTypeId"`
	CiTypeName string                    `json:"ciTypeName"`
	Attributes map[int]map[int]Attribute `json:"attributes"`
}
type Relation struct {
	CiId1            string  `json:"ci_id_1"`
	CiId2            string  `json:"ci_id_2"`
	RelationTypeId   string  `json:"relation_type_id"`
	Direction        string  `json:"direction"`
	RelationTypeName string  `json:"relation_type_name"`
	DirectionName    *string `json:"direction_name"`
}
type Project struct {
	ID                 string `json:"id"`
	Name               string `json:"name"`
	Description        string `json:"description"`
	Note               string `json:"note"`
	OrderNumber        string `json:"order_number"`
	IsActive           string `json:"is_active"`
	UserID             string `json:"user_id"`
	ValidFrom          string `json:"valid_from"`
	CiProjectValidFrom string `json:"ci_project_valid_from"`
	CiProjectHistoryID string `json:"ci_project_history_id"`
}
type Attribute struct {
	ID                   string `json:"id"`
	Name                 string `json:"name"`
	Description          string `json:"description"`
	Note                 string `json:"note"`
	Hint                 string `json:"hint"`
	AttributeTypeID      string `json:"attribute_type_id"`
	AttributeGroupID     string `json:"attribute_group_id"`
	OrderNumber          string `json:"order_number"`
	Column               string `json:"column"`
	IsUnique             string `json:"is_unique"`
	IsNumeric            string `json:"is_numeric"`
	IsBold               string `json:"is_bold"`
	IsEvent              string `json:"is_event"`
	IsUniqueCheck        string `json:"is_unique_check"`
	IsAutocomplete       string `json:"is_autocomplete"`
	IsMultiselect        string `json:"is_multiselect"`
	IsProjectRestricted  string `json:"is_project_restricted"`
	Regex                string `json:"regex"`
	WorkflowID           int    `json:"workflow_id"`
	ScriptName           string `json:"script_name"`
	Tag                  string `json:"tag"`
	InputMaxlength       string `json:"input_maxlength"`
	TextareaCols         string `json:"textarea_cols"`
	TextareaRows         string `json:"textarea_rows"`
	IsActive             string `json:"is_active"`
	UserID               string `json:"user_id"`
	ValidFrom            string `json:"valid_from"`
	Historicize          string `json:"historicize"`
	DisplayStyle         string `json:"display_style"`
	AttributeTypeName    string `json:"attributeTypeName"`
	AttributeGroup       string `json:"attribute_group"`
	ParentAttributeGroup string `json:"parent_attribute_group"`
	ValueText            string `json:"value_text"`
	ValueDate            string `json:"value_date"`
	ValueCi              string `json:"value_ci"`
	CiAttributeID        string `json:"ciAttributeId"`
	Initial              string `json:"initial"`
	ValueNote            string `json:"valueNote"`
	HistoryID            string `json:"history_id"`
	ValueDefault         string `json:"value_default"`
}

func (ciDetail *CiDetail) GetFirstAttributeByName(name string) *Attribute {
	for _, attrMap := range ciDetail.Attributes {
		for _, attr := range attrMap {
			if attr.Name == name {
				return &attr
			}
		}
	}

	return nil
}

func (ciDetail *CiDetail) GetFirstAttributeValueTextByName(name string) (string, bool) {
	attribute := ciDetail.GetFirstAttributeByName(name)
	if attribute == nil {
		return "", false
	}

	return attribute.ValueText, true
}

func (ciDetail *CiDetail) GetFirstAttributeValueCiByName(name string) (int, bool) {
	attribute := ciDetail.GetFirstAttributeByName(name)
	if attribute == nil {
		return 0, false
	}

	ciId, err := strconv.Atoi(attribute.ValueCi)
	if err != nil {
		log.Fatalf("failed to convert ValueCi of \"%s\" to int: %s", name, err.Error())
	}

	return ciId, true
}

type getWorkflowContextResponse struct {
	Status string `json:"status"`
	Data   []struct {
		WorkflowContext string `json:"WorkflowContext"`
	} `json:"data"`
}

func (cmdb *Cmdb) GetWorkflowContext(workflowInstanceId int) (workflowContext *WorkflowContext, err error) {
	if err = cmdb.Login(); err != nil {
		return
	}

	params := map[string]string{
		"argv1": strconv.Itoa(workflowInstanceId),
	}

	getWorkflowContextResponse := getWorkflowContextResponse{}
	err = cmdb.Query("int_getWorkflowContext", &getWorkflowContextResponse, params)
	if err != nil {
		return
	}

	switch len(getWorkflowContextResponse.Data) {
	case 0:
		err = ErrNoResult
		return
	case 1:
		context := getWorkflowContextResponse.Data[0].WorkflowContext

		err = json.Unmarshal([]byte(context), &workflowContext)
		if err != nil {
			err = errors.New(err.Error() + ": " + context)
			return
		}

		return
	default:
		err = ErrTooManyResults
		return
	}
}
