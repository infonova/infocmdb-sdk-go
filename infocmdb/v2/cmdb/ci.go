package cmdb

import (
	"fmt"
	"github.com/infonova/infocmdb-lib-go/core/v2/cmdb/client"
	"gopkg.in/resty.v1"
	"strconv"
)

func (i *Cmdb) CiListByCiTypeID(ciTypeID int, out interface{}) (err error) {
	var respErr client.ResponseError
	resp, err := i.Client.NewRequest().
		SetResult(&out).
		SetError(&respErr).
		SetQueryParams(map[string]string{
			"ciTypeId": fmt.Sprintf("%d", ciTypeID),
		}).
		Get("/apiV2/ci/index")

	if resp != nil && resp.IsError() {
		return respErr
	}

	return
}

type GetCiDetailResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    struct {
		Data struct {
			Ci struct {
				ID        string `json:"id"`
				CiTypeID  string `json:"ci_type_id"`
				Icon      string `json:"icon"`
				HistoryID string `json:"history_id"`
				ValidFrom string `json:"valid_from"`
				CreatedAt string `json:"created_at"`
				UpdatedAt string `json:"updated_at"`
			} `json:"ci"`
			CiType struct {
				ID                      string `json:"id"`
				Name                    string `json:"name"`
				Description             string `json:"description"`
				Note                    string `json:"note"`
				ParentCiTypeID          string `json:"parent_ci_type_id"`
				OrderNumber             string `json:"order_number"`
				CreateButtonDescription string `json:"create_button_description"`
				Icon                    string `json:"icon"`
				Query                   string `json:"query"`
				DefaultProjectID        string `json:"default_project_id"`
				DefaultAttributeID      string `json:"default_attribute_id"`
				DefaultSortAttributeID  string `json:"default_sort_attribute_id"`
				IsDefaultSortAsc        string `json:"is_default_sort_asc"`
				IsCiAttach              string `json:"is_ci_attach"`
				IsAttributeAttach       string `json:"is_attribute_attach"`
				Tag                     string `json:"tag"`
				IsTabEnabled            string `json:"is_tab_enabled"`
				IsEventEnabled          string `json:"is_event_enabled"`
				IsActive                string `json:"is_active"`
				UserID                  string `json:"user_id"`
				ValidFrom               string `json:"valid_from"`
			} `json:"ciType"`
			HistoryCreated string `json:"historyCreated"`
			HistoryChanged string `json:"historyChange"`
			ProjectList    []struct {
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
			} `json:"projectList"`
			AttributeList map[string]struct {
				ID          string `json:"id"`
				Name        string `json:"name"`
				Description string `json:"description"`
				Columns     int    `json:"columns"`
				Attributes  map[string][]struct {
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
					WorkflowID           string `json:"workflow_id"`
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
					ValueText            string `json:"value_text,omitempty"`
					ValueDate            string `json:"value_date"`
					ValueCi              string `json:"value_ci"`
					CiAttributeID        string `json:"ciAttributeId"`
					Initial              string `json:"initial"`
					ValueNote            string `json:"valueNote"`
					HistoryID            string `json:"history_id"`
					ValueDefault         string `json:"value_default"`
					PermissionWrite      string `json:"permission_write"`
				} `json:"attributes"`
				ReadCount  int `json:"readCount"`
				WriteCount int `json:"writeCount"`
			} `json:"attributeList"`
			Icon        string        `json:"icon"`
			Relations   []interface{} `json:"relations"`
			Breadcrumbs []struct {
				ID                      string      `json:"id,omitempty"`
				Name                    string      `json:"name,omitempty"`
				Description             interface{} `json:"description"`
				Note                    string      `json:"note,omitempty"`
				ParentCiTypeID          string      `json:"parent_ci_type_id,omitempty"`
				OrderNumber             string      `json:"order_number,omitempty"`
				CreateButtonDescription interface{} `json:"create_button_description,omitempty"`
				Icon                    interface{} `json:"icon,omitempty"`
				Query                   interface{} `json:"query,omitempty"`
				DefaultProjectID        interface{} `json:"default_project_id,omitempty"`
				DefaultAttributeID      interface{} `json:"default_attribute_id,omitempty"`
				DefaultSortAttributeID  interface{} `json:"default_sort_attribute_id,omitempty"`
				IsDefaultSortAsc        string      `json:"is_default_sort_asc,omitempty"`
				IsCiAttach              string      `json:"is_ci_attach,omitempty"`
				IsAttributeAttach       string      `json:"is_attribute_attach,omitempty"`
				Tag                     interface{} `json:"tag,omitempty"`
				IsTabEnabled            string      `json:"is_tab_enabled,omitempty"`
				IsEventEnabled          string      `json:"is_event_enabled,omitempty"`
				IsActive                string      `json:"is_active,omitempty"`
				UserID                  string      `json:"user_id,omitempty"`
				ValidFrom               string      `json:"valid_from,omitempty"`
				CrumbType               string      `json:"crumbType"`
			} `json:"breadcrumbs"`
			Tickets   []interface{} `json:"tickets"`
			Ticketurl string        `json:"ticketurl"`
			Events    []struct {
				ID                   string      `json:"id"`
				Name                 string      `json:"name"`
				Description          string      `json:"description"`
				Note                 string      `json:"note"`
				Hint                 string      `json:"hint"`
				AttributeTypeID      string      `json:"attribute_type_id"`
				AttributeGroupID     string      `json:"attribute_group_id"`
				OrderNumber          string      `json:"order_number"`
				Column               string      `json:"column"`
				IsUnique             string      `json:"is_unique"`
				IsNumeric            string      `json:"is_numeric"`
				IsBold               string      `json:"is_bold"`
				IsEvent              string      `json:"is_event"`
				IsUniqueCheck        string      `json:"is_unique_check"`
				IsAutocomplete       string      `json:"is_autocomplete"`
				IsMultiselect        string      `json:"is_multiselect"`
				IsProjectRestricted  string      `json:"is_project_restricted"`
				Regex                interface{} `json:"regex"`
				WorkflowID           string      `json:"workflow_id"`
				ScriptName           interface{} `json:"script_name"`
				Tag                  string      `json:"tag"`
				InputMaxlength       interface{} `json:"input_maxlength"`
				TextareaCols         interface{} `json:"textarea_cols"`
				TextareaRows         interface{} `json:"textarea_rows"`
				IsActive             string      `json:"is_active"`
				UserID               string      `json:"user_id"`
				ValidFrom            string      `json:"valid_from"`
				Historicize          string      `json:"historicize"`
				DisplayStyle         interface{} `json:"display_style"`
				AttributeTypeName    string      `json:"attributeTypeName"`
				AttributeGroup       string      `json:"attribute_group"`
				ParentAttributeGroup interface{} `json:"parent_attribute_group"`
				ValueText            interface{} `json:"value_text"`
				ValueDate            interface{} `json:"value_date"`
				ValueCi              interface{} `json:"value_ci"`
				CiAttributeID        string      `json:"ciAttributeId"`
				Initial              string      `json:"initial"`
				ValueNote            interface{} `json:"valueNote"`
				HistoryID            string      `json:"history_id"`
				ValueDefault         interface{} `json:"value_default"`
				PermissionWrite      string      `json:"permission_write"`
			} `json:"events"`
		} `json:"data"`
	} `json:"data"`
}

func (i *Cmdb) CiDetailByCiId(ciId int64) (ciDetail GetCiDetailResponse, restyRes *resty.Response, err error) {
	if err = i.Login(); err != nil {
		return
	}

	var respErr client.ResponseError

	resp, err := i.Client.NewRequest().
		SetResult(&ciDetail).
		SetError(&respErr).
		SetQueryParams(map[string]string{
			"id": strconv.FormatInt(ciId, 10),
		}).Get("/apiV2/ci")

	if resp != nil && resp.IsError() {
		return ciDetail, restyRes, respErr
	}

	return
}
