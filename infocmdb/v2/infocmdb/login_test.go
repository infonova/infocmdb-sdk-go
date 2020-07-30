package infocmdb

import (
	"testing"

	utilTesting "github.com/infonova/infocmdb-sdk-go/util/testing"
)

func TestInfoCMDB_LoginWithUserPass(t *testing.T) {
	url := utilTesting.New().GetUrl()

	type fields struct {
		Config Config
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
			fields{Config{Url: url, Username: "admin", Password: "admin"}},
			args{},
			false,
		},
		{
			"invalid login admin//noadmin",
			fields{Config{Url: url, Username: "admin", Password: "fail"}},
			args{},
			true,
		},
		{
			"invalid no username",
			fields{Config{Url: url, Username: "", Password: "fail"}},
			args{},
			true,
		},
		{
			"invalid no password",
			fields{Config: Config{Url: url, Username: "admin", Password: ""}},
			args{},
			true,
		},
		{
			"invalid no data",
			fields{Config: Config{Url: url, Username: "", Password: ""}},
			args{},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmdbV2 := New()
			cmdbV2.LoadConfig(tt.fields.Config)

			if err := cmdbV2.Login(); (err != nil) != tt.wantErr {
				t.Errorf("Cmdb.LoginWithUserPass() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
