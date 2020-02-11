package infocmdb

import (
	"reflect"
	"testing"

	v2 "github.com/infonova/infocmdb-sdk-go/infocmdb/v2/infocmdb"
	utilTesting "github.com/infonova/infocmdb-sdk-go/util/testing"
)

func TestInfoCMDB_GetListOfCiIdsOfCiType(t *testing.T) {
	infocmdbUrl := utilTesting.New().GetUrl()
	type fields struct {
		v2Config v2.Config
	}
	type args struct {
		ciTypeID int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantR   CiIds
		wantErr bool
	}{
		{
			"v2 List Ci's pf Type '1' with wrong Credentials (fail)",
			fields{
				v2.Config{
					Url:      infocmdbUrl,
					Username: "fail",
					Password: "fail",
					BasePath: "/app/",
				},
			},
			args{ciTypeID: 1},
			nil,
			true,
		},
		{
			"v2 List Ci's of Type '1' (demo)",
			fields{
				v2.Config{
					Url:      infocmdbUrl,
					Username: "admin",
					Password: "admin",
					BasePath: "/app/",
				},
			},
			args{ciTypeID: 1},
			CiIds{1, 2},
			false,
		},
		{
			"v2 List Ci's of Type '-1' (error)",
			fields{
				v2.Config{
					Url:      infocmdbUrl,
					Username: "admin",
					Password: "admin",
					BasePath: "/app/",
				},
			},
			args{ciTypeID: -1},
			nil,
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
			gotR, err := cmdb.GetListOfCiIdsOfCiType(tt.args.ciTypeID)
			if (err != nil) != tt.wantErr {
				t.Errorf("getListOfCiIdsOfCiType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotR, tt.wantR) {
				t.Errorf("getListOfCiIdsOfCiType() gotR = %v, want %v", gotR, tt.wantR)
			}
		})
	}
}

func TestInfoCMDB_GetListOfCiIdsOfCiTypeV2(t *testing.T) {
	infocmdbUrl := utilTesting.New().GetUrl()
	type fields struct {
		v2Config v2.Config
	}
	type args struct {
		ciTypeID int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantR   CiIds
		wantErr bool
	}{
		{
			"v2 List Ci's pf Type '1' with wrong Credentials (fail)",
			fields{
				v2.Config{
					Url:      infocmdbUrl,
					Username: "false",
					Password: "false",
					BasePath: "/app/",
				},
			},
			args{ciTypeID: 1},
			nil,
			true,
		},
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
			args{ciTypeID: 1},
			CiIds{1, 2},
			false,
		},
		{
			"v2 List Ci's of Type '-1' (error)",
			fields{
				v2.Config{
					Url:      infocmdbUrl,
					Username: "admin",
					Password: "admin",
					BasePath: "/app/",
				},
			},
			args{ciTypeID: -1},
			nil,
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
			gotR, err := cmdb.GetListOfCiIdsOfCiTypeV2(tt.args.ciTypeID)
			if (err != nil) != tt.wantErr {
				t.Errorf("getListOfCiIdsOfCiType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotR, tt.wantR) {
				t.Errorf("getListOfCiIdsOfCiType() gotR = %v, want %v", gotR, tt.wantR)
			}
		})
	}
}

func TestInfoCMDB_CreateCi(t *testing.T) {
	infocmdbUrl := utilTesting.New().GetUrl()
	type fields struct {
		v2Config v2.Config
	}
	type args struct {
		ciTypeID  int
		icon      string
		historyId int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantR   CreateCi
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
				ciTypeID:  476,
				icon:      "",
				historyId: 0,
			},
			CreateCi{
				ID:        617827,
				CiTypeID:  476,
				Icon:      "",
				HistoryID: 59529024,
				ValidFrom: "2020-01-13 15:14:05",
				CreatedAt: "2020-01-13 15:14:05",
				UpdatedAt: "2020-01-13 15:14:05",
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
			gotR, err := cmdb.CreateCi(tt.args.ciTypeID, tt.args.icon, tt.args.historyId)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateCi() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotR, tt.wantR) {
				t.Errorf("CreateCi() gotR = %v, want %v", gotR, tt.wantR)
			}
		})
	}
}

