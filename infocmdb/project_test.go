package infocmdb

import (
	"testing"

	v2 "github.com/infonova/infocmdb-sdk-go/infocmdb/v2/infocmdb"
	utilTesting "github.com/infonova/infocmdb-sdk-go/util/testing"
)

func TestInfoCMDB_AddCiProjectMapping(t *testing.T) {
	infocmdbUrl := utilTesting.New().GetUrl()
	type fields struct {
		v2Config v2.Config
	}
	type args struct {
		ciID      int
		projectID int
		historyID int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"v2 create CI",
			fields{
				v2.Config{
					Url:      infocmdbUrl,
					Username: "admin",
					Password: "admin",
					BasePath: "/app/",
				},
			},
			args{
				ciID:      617830,
				projectID: 33,
				historyID: 59529029,
			},
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
			err := cmdb.AddCiProjectMapping(tt.args.ciID, tt.args.projectID, tt.args.historyID)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddCiProjectMapping() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
