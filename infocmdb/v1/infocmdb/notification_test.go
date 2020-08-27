package infocmdb

import (
	"reflect"
	"testing"

	utilTesting "github.com/infonova/infocmdb-sdk-go/util/testing"
)

func TestInfoCMDB_SendNotification(t *testing.T) {
	ut := utilTesting.New().AddMocking(utilTesting.Mocking{
		RequestString: `POST##/api/notification/notify/default/method/json##apikey=mocked_mocked_mocked_mocked%23%23%23&one=two&three=four`,
		ReturnString:  `{"status":"OK","data":[{"type":"mail","address":"infocmdb@localhost"},{"type":"mail","address":"marko.oslaj@bearingpoint.com"}]}`,
	}).AddMocking(utilTesting.Mocking{
		RequestString: `POST##/api/notification/notify/bla/method/json##apikey=mocked_mocked_mocked_mocked%23%23%23`,
		ReturnString:  `{"status":"error","message":"unexpected Error occurred."}`,
	})

	type args struct {
		notifyName string
		options    NotifyParams
	}

	mapica := make(map[string]string)

	mapica["one"] = "two"
	mapica["three"] = "four"

	var tests = []struct {
		name    string
		args    args
		wantR   NotificationResponse
		wantErr bool
	}{
		{
			"v1 SendNotification",
			args{
				notifyName: "default",
				options: NotifyParams{
					Recipients:  []string{"marko.oslaj@bearingpoint.com"},
					Subject:     "Test",
					OtherParams: mapica,
				},
			},
			NotificationResponse{SentTo: []SendNotificationData{
				{
					Type:    "mail",
					Address: "infocmdb@localhost",
				},
				{
					Type:    "mail",
					Address: "marko.oslaj@bearingpoint.com",
				},
			}},
			false,
		},
		{
			"v1 SendNotification 2",
			args{
				notifyName: "bla",
				options:    NotifyParams{},
			},
			NotificationResponse{},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmdb := New()
			ut.SetValidConfig(&cmdb.Config)
			cmdb.LoadConfig(cmdb.Config)

			if err := cmdb.Login(); err != nil {
				t.Fatalf("Login failed: %v\n", err)
			}

			gotR, err := cmdb.SendNotification(tt.args.notifyName, tt.args.options)
			if (err != nil) != tt.wantErr {
				t.Errorf("SendNotification() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(gotR, tt.wantR) {
				t.Errorf("SendNotification() gotR = %v, want %v", gotR, tt.wantR)
			}

		})
	}
}
