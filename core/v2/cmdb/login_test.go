package cmdb

import (
	"os"
	"testing"

	"github.com/infonova/infocmdb-lib-go/core/v2/cmdb/client"
	"github.com/patrickmn/go-cache"
)

func getTestUrl() (url string) {
	if url = os.Getenv("WORKFLOW_TEST_URL"); url == "" {
		url = "http://localhost"
	}

	return
}

func TestInfoCMDB_LoginWithUserPass(t *testing.T) {
	type fields struct {
		Config Config
		Cache  *cache.Cache
		Client *client.Client
	}
	type args struct {
		url      string
		username string
		password string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"valid login admin//admin",
			fields{Config: Config{Url: getTestUrl(), Username: "admin", Password: "admin"}},
			args{},
			false,
		},
		{
			"invalid login admin//noadmin",
			fields{Config: Config{Url: getTestUrl(), Username: "admin", Password: "noadmin"}},
			args{},
			true,
		},
		{
			"invalid no username",
			fields{Config: Config{Url: getTestUrl(), Username: "", Password: "noadmin"}},
			args{},
			true,
		},
		{
			"invalid no password",
			fields{Config: Config{Url: getTestUrl(), Username: "admin", Password: ""}},
			args{},
			true,
		},
		{
			"invalid no data",
			fields{Config: Config{Url: "", Username: "", Password: ""}},
			args{},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := new(InfoCMDB)
			i.Config = tt.fields.Config
			i.Cache = tt.fields.Cache
			i.Client = client.NewClient(i.Config.Url)

			if err := i.Login(); (err != nil) != tt.wantErr {
				t.Errorf("InfoCMDB.LoginWithUserPass() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
