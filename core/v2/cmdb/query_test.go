package cmdb

import (
	"os"
	"testing"

	"github.com/infonova/infocmdb-lib-go/core/v2/cmdb/client"
	"github.com/patrickmn/go-cache"
)

func init() {
	if infocmdbUrl = os.Getenv("WORKFLOW_TEST_URL"); infocmdbUrl == "" {
		infocmdbUrl = "http://localhost"
	}
}

func TestInfoCMDB_Query(t *testing.T) {
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
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"valid", fields{Config: Config{Url: infocmdbUrl, Username: "admin", Password: "admin"}}, args{}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &InfoCMDB{
				Config: tt.fields.Config,
				Cache:  tt.fields.Cache,
				Client: client.NewClient(tt.fields.Config.Url),
				Error:  tt.fields.Error,
			}
			if err := i.Query(tt.args.query, tt.args.out, tt.args.params); (err != nil) != tt.wantErr {
				t.Errorf("InfoCMDB.Query() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
