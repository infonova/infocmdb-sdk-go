package core

import (
	v1 "github.com/infonova/infocmdb-lib-go/core/v1/cmdb"
	v2 "github.com/infonova/infocmdb-lib-go/core/v2/cmdb"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"os"
	"reflect"
	"testing"
)

var (
	mocking            = false
	infoCMDBConfigFile = "test/test.yml"
	infocmdbUrl        = "http://localhost"
)

func init() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("failed to load env: %v", err)
	}

	if os.Getenv("WORKFLOW_TEST_MOCKING") == "true" {
		mocking = true
		log.Debug("Mocking enabled")
	}

	infocmdbUrl = os.Getenv("WORKFLOW_TEST_URL")
}

func TestInfoCMDB_GetListOfCiIdsOfCiType(t *testing.T) {
	type fields struct {
		v1 *v1.InfoCMDB
		v2 *v2.InfoCMDB
	}
	type args struct {
		ciTypeID int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantR   ListOfCiIdsOfCiType
		wantErr bool
	}{
		{
			"v1 List Ci's pf Type '1' with wrong Credentials (fail)",
			fields{
				&v1.InfoCMDB{Config: v1.Config{
					ApiUrl:       infocmdbUrl,
					ApiUser:      "fail",
					ApiPassword:  "fail",
					CmdbBasePath: "/app/",
				}},
				&v2.InfoCMDB{Config: v2.Config{}},
			},
			args{ciTypeID: 1},
			nil,
			true,
		},
		{
			"v1 List Ci's of Type '1' (demo)",
			fields{
				&v1.InfoCMDB{Config: v1.Config{
					ApiUrl:       infocmdbUrl,
					ApiUser:      "admin",
					ApiPassword:  "admin",
					CmdbBasePath: "/app/",
				}},
				&v2.InfoCMDB{Config: v2.Config{}},
			},
			args{ciTypeID: 1},
			ListOfCiIdsOfCiType{{1}, {2}},
			false,
		},
		{
			"v1 List Ci's of Type '-1' (error)",
			fields{
				&v1.InfoCMDB{Config: v1.Config{
					ApiUrl:       infocmdbUrl,
					ApiUser:      "admin",
					ApiPassword:  "admin",
					CmdbBasePath: "/app/",
				}},
				&v2.InfoCMDB{Config: v2.Config{}},
			},
			args{ciTypeID: -1},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &InfoCMDB{
				v1: tt.fields.v1,
				v2: tt.fields.v2,
			}
			gotR, err := i.GetListOfCiIdsOfCiType(tt.args.ciTypeID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetListOfCiIdsOfCiType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotR, tt.wantR) {
				t.Errorf("GetListOfCiIdsOfCiType() gotR = %v, want %v", gotR, tt.wantR)
			}
		})
	}
}

func TestInfoCMDB_GetListOfCiIdsOfCiTypeV2(t *testing.T) {
	type fields struct {
		v1 *v1.InfoCMDB
		v2 *v2.InfoCMDB
	}
	type args struct {
		ciTypeID int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantR   ListOfCiIdsOfCiType
		wantErr bool
	}{
		{
			"v2 List Ci's pf Type '1' with wrong Credentials (fail)",
			fields{
				&v1.InfoCMDB{Config: v1.Config{}},
				&v2.InfoCMDB{Config: v2.Config{
					Url:      infocmdbUrl,
					Username: "false",
					Password: "false",
					BasePath: "/app/",
				}},
			},
			args{ciTypeID: 1},
			nil,
			true,
		},
		{
			"v2 List CIs of Type 1 (demo)",
			fields{
				&v1.InfoCMDB{Config: v1.Config{}},
				&v2.InfoCMDB{Config: v2.Config{
					Url:      infocmdbUrl,
					Username: "admin",
					Password: "admin",
					BasePath: "/app/",
				}},
			},
			args{ciTypeID: 1},
			ListOfCiIdsOfCiType{{1}, {2}},
			false,
		},
		{
			"v2 List Ci's of Type '-1' (error)",
			fields{
				&v1.InfoCMDB{Config: v1.Config{}},
				&v2.InfoCMDB{Config: v2.Config{
					Url:      infocmdbUrl,
					Username: "admin",
					Password: "admin",
					BasePath: "/app/",
				}},
			},
			args{ciTypeID: -1},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &InfoCMDB{
				v1: tt.fields.v1,
				v2: tt.fields.v2,
			}
			gotR, err := i.GetListOfCiIdsOfCiTypeV2(tt.args.ciTypeID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetListOfCiIdsOfCiType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotR, tt.wantR) {
				t.Errorf("GetListOfCiIdsOfCiType() gotR = %v, want %v", gotR, tt.wantR)
			}
		})
	}
}

func TestInfoCMDB_QueryWebservice(t *testing.T) {
	type fields struct {
		v2 *v2.InfoCMDB
	}
	type args struct {
		ws     string
		params map[string]string
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantR   string
		wantErr bool
	}{
		{
			"v2 List CIs of Type 1 (demo)",
			fields{
				&v2.InfoCMDB{Config: v2.Config{
					Url:      infocmdbUrl,
					Username: "admin",
					Password: "admin",
					BasePath: "/app/",
				}},
			},
			args{ws: "int_getCi", params: map[string]string{"argv1": "1"}},
			`{"success":true,"message":"Query executed successfully","data":[{"ci_id":"1","ci_type_id":"1","ci_type":"demo","project":"springfield","project_id":"4"}]}`,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &InfoCMDB{
				v2: tt.fields.v2,
			}
			gotR, err := i.QueryWebservice(tt.args.ws, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("QueryWebservice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotR != tt.wantR {
				t.Errorf("QueryWebservice() gotR = %v, want %v", gotR, tt.wantR)
			}
		})
	}
}

func TestInfoCMDB_UpdateCiAttribute(t *testing.T) {
	type fields struct {
		v2 *v2.InfoCMDB
	}
	type args struct {
		ci int
		ua []UpdateCiAttribute
	}

	cmdbConfigValid := v2.InfoCMDB{Config: v2.Config{
		Url:      infocmdbUrl,
		Username: "admin",
		Password: "admin",
		BasePath: "/app/",
	}}
	baseCiID := 14 // ci to base this tests on
	baseAttributeName := "emp_lastname"

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"v2 Delete CI Attribute - fail - requires ciattributeid",
			fields{&cmdbConfigValid},
			args{ci: baseCiID, ua: []UpdateCiAttribute{
				{Mode: v2.UpdateModeDelete, Name: baseAttributeName},
			}},
			true,
		},
		{
			"v2 Update CI Attribute",
			fields{&cmdbConfigValid},
			args{ci: baseCiID, ua: []UpdateCiAttribute{
				{Mode: v2.UpdateModeSet, Name: "emp_firstname", Value: "22322"},
			}},
			false,
		},
		{
			"v2 Update CI Attribute - wrong attribute name",
			fields{&cmdbConfigValid},
			args{ci: baseCiID, ua: []UpdateCiAttribute{
				{Mode: v2.UpdateModeSet, Name: baseAttributeName + "_NOT_EXISTING", Value: "1"},
			}},
			true,
		},
		{
			"v2 Insert CI Attribute",
			fields{&cmdbConfigValid},
			args{ci: baseCiID, ua: []UpdateCiAttribute{
				{Mode: v2.UpdateModeInsert, Name: baseAttributeName, Value: "New1"},
			}},
			false,
		},
		{
			"v2 Update CI Attribute - fail - multiple attributes with the same name",
			fields{&cmdbConfigValid},
			args{ci: baseCiID, ua: []UpdateCiAttribute{
				{Mode: v2.UpdateModeSet, Name: baseAttributeName, Value: "22322"},
			}},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &InfoCMDB{
				v2: tt.fields.v2,
			}
			err := i.UpdateCiAttribute(tt.args.ci, tt.args.ua)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateCiAttribute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
