package infocmdb

import (
	"testing"

	v2 "github.com/infonova/infocmdb-sdk-go/infocmdb/v2/infocmdb"
	utilTesting "github.com/infonova/infocmdb-sdk-go/util/testing"
)

func TestInfoCMDB_GetAttrDefaultOptionIdByAttrId(t *testing.T) {
	infocmdbUrl := utilTesting.New().GetUrl()
	type fields struct {
		v2Config v2.Config
	}
	type args struct {
		attrId      int
		optionValue string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantR   int
		wantErr bool
	}{
		{
			"v2 GetAttributeDefaultOptionId",
			fields{
				v2.Config{
					Url:      infocmdbUrl,
					Username: "admin",
					Password: "admin",
					BasePath: "/app/",
				},
			},
			args{
				attrId:      438,
				optionValue: "IN PROGRESS",
			},
			1329,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmdbV2 := v2.New()
			if err := cmdbV2.LoadConfig(tt.fields.v2Config); err != nil {
				t.Fatalf("LoadConfig failed: %v\n", err)
			}
			cmdb := &Client{
				v2: cmdbV2,
			}
			gotR, err := cmdb.GetAttrDefaultOptionIdByAttrId(tt.args.attrId, tt.args.optionValue)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAttrDefaultOptionIdByAttrId() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotR != tt.wantR {
				t.Errorf("GetAttrDefaultOptionIdByAttrId() gotR = %v, want %v", gotR, tt.wantR)
			}
		})
	}
}

func TestInfoCMDB_UpdateCiAttribute(t *testing.T) {
	infocmdbUrl := utilTesting.New().GetUrl()
	type fields struct {
		v2Config v2.Config
	}
	type args struct {
		ci int
		ua []v2.UpdateCiAttribute
	}

	cmdbConfigValid := v2.Config{
		Url:      infocmdbUrl,
		Username: "admin",
		Password: "admin",
		BasePath: "/app/",
	}
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
			fields{cmdbConfigValid},
			args{ci: baseCiID, ua: []v2.UpdateCiAttribute{
				{Mode: v2.UPDATE_MODE_DELETE, Name: baseAttributeName},
			}},
			true,
		},
		{
			"v2 Update CI Attribute",
			fields{cmdbConfigValid},
			args{ci: baseCiID, ua: []v2.UpdateCiAttribute{
				{Mode: v2.UPDATE_MODE_SET, Name: "emp_firstname", Value: "22322"},
			}},
			false,
		},
		{
			"v2 Update CI Attribute - wrong attribute name",
			fields{cmdbConfigValid},
			args{ci: baseCiID, ua: []v2.UpdateCiAttribute{
				{Mode: v2.UPDATE_MODE_SET, Name: baseAttributeName + "_NOT_EXISTING", Value: "1"},
			}},
			true,
		},
		{
			"v2 Insert CI Attribute",
			fields{cmdbConfigValid},
			args{ci: baseCiID, ua: []v2.UpdateCiAttribute{
				{Mode: v2.UPDATE_MODE_INSERT, Name: baseAttributeName, Value: "New1"},
			}},
			false,
		},
		{
			"v2 Update CI Attribute - fail - multiple attributes with the same name",
			fields{cmdbConfigValid},
			args{ci: baseCiID, ua: []v2.UpdateCiAttribute{
				{Mode: v2.UPDATE_MODE_SET, Name: baseAttributeName, Value: "22322"},
			}},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmdbV2 := v2.New()
			if err := cmdbV2.LoadConfig(tt.fields.v2Config); err != nil {
				t.Fatalf("LoadConfig failed: %v\n", err)
			}
			cmdb := &Client{
				v2: cmdbV2,
			}
			err := cmdb.UpdateCiAttribute(tt.args.ci, tt.args.ua)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateCiAttribute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

