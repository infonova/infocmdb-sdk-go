package infocmdb

import (
	"testing"

	utilTesting "github.com/infonova/infocmdb-sdk-go/util/testing"
)

func TestInfoCMDB_Query(t *testing.T) {
	url := utilTesting.New().GetUrl()

	type fields struct {
		Config Config
	}
	type args struct {
		query  string
		out    interface{}
		params map[string]string
	}
	var out interface{}
	var p = map[string]string{
		"argv1": "428",
		"argv2": "29",
	}
	a := args{
		query:  "int_getCiAttributeId",
		out:    &out,
		params: p,
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"valid", fields{Config{Url: url, Username: "admin", Password: "admin"}}, a, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmdbV2 := New()
			cmdbV2.LoadConfig(tt.fields.Config)

			if err := cmdbV2.Query(tt.args.query, tt.args.out, tt.args.params); (err != nil) != tt.wantErr {
				t.Errorf("Cmdb.Query() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
