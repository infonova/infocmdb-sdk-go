package infocmdb

import (
	"testing"

	v2 "github.com/infonova/infocmdb-sdk-go/infocmdb/v2/infocmdb"
	utilTesting "github.com/infonova/infocmdb-sdk-go/util/testing"
)

func TestInfoCMDB_QueryWebservice(t *testing.T) {
	infocmdbUrl := utilTesting.New().GetUrl()
	type fields struct {
		v2Config v2.Config
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
				v2.Config{
					Url:      infocmdbUrl,
					Username: "admin",
					Password: "admin",
					BasePath: "/app/",
				},
			},
			args{ws: "int_getCi", params: map[string]string{"argv1": "1"}},
			`{"success":true,"message":"Query executed successfully","data":[{"ci_id":"1","ci_type_id":"1","ci_type":"demo","project":"springfield","project_id":"4"}]}`,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmdbV2 := v2.New()
			cmdbV2.LoadConfig(tt.fields.v2Config)
			cmdb := &Client{
				v2: cmdbV2,
			}
			gotR, err := cmdb.QueryWebservice(tt.args.ws, tt.args.params)
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
