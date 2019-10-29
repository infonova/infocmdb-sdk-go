package cmdb

import (
	utilTesting "github.com/infonova/infocmdb-lib-go/util/testing"
	"os"
	"testing"
)

func TestInfoCMDB_CiListByCiTypeID(t *testing.T) {
	utilTesting.LoadEnvFromFile("../../.env")
	url := os.Getenv("WORKFLOW_TEST_URL")

	configFile := utilTesting.BuildValidConfig(url)

	type args struct {
		ciTypeID int
		out      interface{}
	}

	type EmployeeReturn struct {
		Data struct {
			Data struct {
				AttributeList map[int]struct {
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
	var citypetest1 EmployeeReturn

	tests := []struct {
		name    string
		Config  Config
		args    args
		wantErr bool
	}{
		{
			"query citype 1",
			Config{},
			args{12, citypetest1},
			false,
		},
		{
			"fail query citype -1",
			Config{},
			args{-1, citypetest1},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i, e := NewCMDB().LoadConfig(configFile)
			if e != nil {
				panic(e)
			}
			if err := i.Login(); err != nil {
				t.Logf("Login failed: %v\n", err)
			}

			var citypetest2 EmployeeReturn
			if err := i.CiListByCiTypeID(tt.args.ciTypeID, &citypetest2); (err != nil) != tt.wantErr {
				t.Errorf("InfoCMDB.CiListByCiTypeID() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
