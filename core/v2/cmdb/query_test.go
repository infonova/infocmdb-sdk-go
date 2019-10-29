package cmdb

import (
	"testing"

	"github.com/infonova/infocmdb-lib-go/core/v2/cmdb/client"
	utilTesting "github.com/infonova/infocmdb-lib-go/util/testing"
	"github.com/patrickmn/go-cache"
)

func TestInfoCMDB_Query(t *testing.T) {
	url := utilTesting.Url

	type fields struct {
		Config Config
		Cache  *cache.Cache
		Client *client.Client
		Error  error
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
		{"valid", fields{Config: Config{Url: url, Username: "admin", Password: "admin"}}, a, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Cmdb{
				Config: tt.fields.Config,
				Cache:  tt.fields.Cache,
				Client: client.NewClient(tt.fields.Config.Url),
				Error:  tt.fields.Error,
			}
			if err := i.Query(tt.args.query, tt.args.out, tt.args.params); (err != nil) != tt.wantErr {
				t.Errorf("Cmdb.Query() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
