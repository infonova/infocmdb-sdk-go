package cmdb

import (
	"fmt"
	"strconv"
)

type EmployeeReturn struct {
	Data struct {
		Data struct {
			AttributeList struct {
				One struct {
					AttributeGroupID    string      `json:"attribute_group_id"`
					AttributeTypeID     string      `json:"attribute_type_id"`
					Column              string      `json:"column"`
					Description         string      `json:"description"`
					DisplayStyle        interface{} `json:"display_style"`
					Hint                string      `json:"hint"`
					Historicize         string      `json:"historicize"`
					ID                  string      `json:"id"`
					InputMaxlength      interface{} `json:"input_maxlength"`
					IsActive            string      `json:"is_active"`
					IsAutocomplete      string      `json:"is_autocomplete"`
					IsBold              string      `json:"is_bold"`
					IsEvent             string      `json:"is_event"`
					IsMultiselect       string      `json:"is_multiselect"`
					IsNumeric           string      `json:"is_numeric"`
					IsProjectRestricted string      `json:"is_project_restricted"`
					IsUnique            string      `json:"is_unique"`
					IsUniqueCheck       string      `json:"is_unique_check"`
					Name                string      `json:"name"`
					Note                string      `json:"note"`
					OrderNumber         string      `json:"order_number"`
					Regex               interface{} `json:"regex"`
					ScriptName          interface{} `json:"script_name"`
					Tag                 string      `json:"tag"`
					TextareaCols        interface{} `json:"textarea_cols"`
					TextareaRows        interface{} `json:"textarea_rows"`
					UserID              string      `json:"user_id"`
					ValidFrom           string      `json:"valid_from"`
					Width               interface{} `json:"width"`
					WorkflowID          interface{} `json:"workflow_id"`
				} `json:"1"`
			} `json:"attributeList"`
			Breadcrumbs []struct {
				CreateButtonDescription interface{} `json:"create_button_description"`
				CrumbType               string      `json:"crumbType"`
				DefaultAttributeID      string      `json:"default_attribute_id"`
				DefaultProjectID        string      `json:"default_project_id"`
				DefaultSortAttributeID  string      `json:"default_sort_attribute_id"`
				Description             string      `json:"description"`
				Icon                    interface{} `json:"icon"`
				ID                      string      `json:"id"`
				IsActive                string      `json:"is_active"`
				IsAttributeAttach       string      `json:"is_attribute_attach"`
				IsCiAttach              string      `json:"is_ci_attach"`
				IsDefaultSortAsc        string      `json:"is_default_sort_asc"`
				IsEventEnabled          string      `json:"is_event_enabled"`
				IsTabEnabled            string      `json:"is_tab_enabled"`
				Name                    string      `json:"name"`
				Note                    string      `json:"note"`
				OrderNumber             string      `json:"order_number"`
				ParentCiTypeID          string      `json:"parent_ci_type_id"`
				Query                   interface{} `json:"query"`
				Tag                     interface{} `json:"tag"`
				UserID                  string      `json:"user_id"`
				ValidFrom               string      `json:"valid_from"`
			} `json:"breadcrumbs"`
			CiList []struct {
				Color             interface{} `json:"color"`
				EmpEmailAddress   string      `json:"emp_email_address"`
				EmpEmailAddressID string      `json:"emp_email_addressID"`
				EmpFirstname      string      `json:"emp_firstname"`
				EmpFirstnameID    string      `json:"emp_firstnameID"`
				EmpLastname       string      `json:"emp_lastname"`
				EmpLastnameID     string      `json:"emp_lastnameID"`
				EmpStaffNumber    string      `json:"emp_staff_number"`
				EmpStaffNumberID  string      `json:"emp_staff_numberID"`
				ID                string      `json:"id"`
			} `json:"ciList"`
			CiTypeAttach            string      `json:"ciTypeAttach"`
			CreateButtonDescription interface{} `json:"createButtonDescription"`
			DefaultOrderBy          interface{} `json:"defaultOrderBy"`
			IsQuery                 bool        `json:"isQuery"`
			ListEdit                string      `json:"listEdit"`
			Paginator               struct{}    `json:"paginator"`
			Scrollbar               bool        `json:"scrollbar"`
			TypeName                string      `json:"typeName"`
		} `json:"data"`
	} `json:"data"`
	Message string `json:"message"`
	Success bool   `json:"success"`
}

func (i *InfoCMDB) CiListByCiTypeID(ciTypeID int, out interface{}) (err error) {
	params := map[string]string{
		"ciTypeId":             fmt.Sprintf("%d", ciTypeID),
		"XDEBUG_SESSION":       "XDEBUG_ECLIPSE",
		"XDEBUG_SESSION_START": "XDEBUG_ECLIPSE",
	}
	err = i.Client.Get("/apiV2/ci/index", &out, params)

	fmt.Printf("Message: %v\n", out)
	if err != nil {
		i.AddError(err)
		return
	}

	return
}

type GetCiResponse struct {
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

func (i *InfoCMDB) CiDetailByCiId(ciId int64) (resp GetCiResponse, err error) {
	params := map[string]string{
		"id": strconv.FormatInt(ciId, 10),
	}

	err = i.Client.Get("/apiV2/ci", &resp, params)

	return
}
