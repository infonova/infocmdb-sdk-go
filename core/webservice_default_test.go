package cmdb

import (
	v1 "git.appteam.infonova.cloud/infocmdb/library/core/v1/cmdb"
	v2 "git.appteam.infonova.cloud/infocmdb/library/core/v2/cmdb"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"os"
	"reflect"
	"testing"
)

var (
	mocking            = false
	infoCMDBConfigFile = "test/test.yml"
	infocmdbUrl = "http://localhost"
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
					ApiUrl: infocmdbUrl,
					ApiUser: "fail",
					ApiPassword: "fail",
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
					ApiUrl: infocmdbUrl,
					ApiUser: "admin",
					ApiPassword: "admin",
					CmdbBasePath: "/app/",
				}},
				&v2.InfoCMDB{Config: v2.Config{}},
			},
			args{ciTypeID: 1},
			ListOfCiIdsOfCiType{{1},{2}},
			false,
		},
		{
			"v1 List Ci's of Type '-1' (error)",
			fields{
				&v1.InfoCMDB{Config: v1.Config{
					ApiUrl: infocmdbUrl,
					ApiUser: "admin",
					ApiPassword: "admin",
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
			"v2 List Ci's of Type '1' (demo)",
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
			ListOfCiIdsOfCiType{{1},{2}},
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