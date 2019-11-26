package infocmdb

import (
	"bytes"
	"fmt"
	v1 "github.com/infonova/infocmdb-sdk-go/infocmdb/v1/infocmdb"
	v2 "github.com/infonova/infocmdb-sdk-go/infocmdb/v2/infocmdb"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"reflect"
	"testing"
)

var (
	mocking     bool
	infocmdbUrl string
)

func init() {
	if os.Getenv("WORKFLOW_TEST_MOCKING") == "true" {
		mocking = true
		log.Debug("Mocking enabled")
	} else if err := godotenv.Load("../../.env"); err != nil {
		log.Fatalf("failed to load env: %v", err)
	}

	err := godotenv.Load("../../.env")

	if err != nil {
		if mocking {
			log.Infof("ignoring failure to load env due to enabled mocking, error: %v", err)
		} else {
			log.Fatalf("failed to load env: %v", err)
		}
	}

	infocmdbUrl = getUrl()
}

func newMockServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		infocmdbUrl := os.Getenv("WORKFLOW_TEST_URL") + r.URL.String()
		mockString := fmt.Sprintf("%s##%s##%s", r.Method, r.URL.String(), string(body))
		switch mockString {
		case `POST##/apiV2/auth/token##lifetime=600&password=fail&username=fail`:
			w.Header()["Content-Type"] = []string{"application/json;charset=UTF-8"}
			w.WriteHeader(http.StatusUnauthorized)
			return
		case `POST##/apiV2/auth/token##lifetime=600&password=false&username=false`:
			w.Header()["Content-Type"] = []string{"application/json;charset=UTF-8"}
			w.WriteHeader(http.StatusBadRequest)
			return
		case `POST##/apiV2/auth/token##lifetime=600&password=admin&username=admin`:
			w.Header()["Content-Type"] = []string{"application/json;charset=UTF-8"}
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"success":true,"message":"mocked","data":{"token":"mocked"}}`))
			return
		case `PUT##/apiV2/query/execute/int_getListOfCiIdsOfCiType##{"query":{"params":{"argv1":"1"}}}`:
			w.Header()["Content-Type"] = []string{"application/json;charset=UTF-8"}
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"status":"OK","data":[{"ciid":"1"},{"ciid":"2"}]}`))
			return
		case `PUT##/apiV2/query/execute/int_getCi##{"query":{"params":{"argv1":"1"}}}`:
			w.Header()["Content-Type"] = []string{"application/json;charset=UTF-8"}
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"success":true,"message":"Query executed successfully","data":[{"ci_id":"1","ci_type_id":"1","ci_type":"demo","project":"springfield","project_id":"4"}]}`))
			return
		case `PUT##/apiV2/ci/14##{"ci":{"attributes":[{"mode":"set","name":"emp_firstname","value":"22322","ciAttributeId":0,"uploadId":""}]}}`:
			w.Header()["Content-Type"] = []string{"application/json;charset=UTF-8"}
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"success":true,"message":"Query executed successfully","data":[]}`))
			return
		case `PUT##/apiV2/ci/14##{"ci":{"attributes":[{"mode":"insert","name":"emp_lastname","value":"New1","ciAttributeId":0,"uploadId":""}]}}`:
			w.Header()["Content-Type"] = []string{"application/json;charset=UTF-8"}
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"success":true,"message":"Query executed successfully","data":[]}`))
			return
		case `PUT##/apiV2/ci/14##{"ci":{"attributes":[{"mode":"delete","name":"emp_lastname","value":"","ciAttributeId":0,"uploadId":""}]}}`:
			w.Header()["Content-Type"] = []string{"application/json;charset=UTF-8"}
			w.WriteHeader(http.StatusBadRequest)
			return
		case `PUT##/apiV2/ci/14##{"ci":{"attributes":[{"mode":"set","name":"emp_lastname_NOT_EXISTING","value":"1","ciAttributeId":0,"uploadId":""}]}}`:
			w.Header()["Content-Type"] = []string{"application/json;charset=UTF-8"}
			w.WriteHeader(http.StatusBadRequest)
			return
		case `PUT##/apiV2/ci/14##{"ci":{"attributes":[{"mode":"set","name":"emp_lastname","value":"22322","ciAttributeId":0,"uploadId":""}]}}`:
			w.Header()["Content-Type"] = []string{"application/json;charset=UTF-8"}
			w.WriteHeader(http.StatusBadRequest)
			return
		default:
			fmt.Printf("Didn't Mock:\n''''''\n%s\n''''''\n", mockString)
		}
		// you can reassign the body if you need to parse it as multipart
		r.Body = ioutil.NopCloser(bytes.NewReader(body))
		body, _ = ioutil.ReadAll(r.Body)
		values, _ := url.ParseQuery(string(body))

		fmt.Printf("Method: %s, Url: %s, Data: %s\n", r.Method, infocmdbUrl, body)
		fmt.Printf("%v\n", values)

		// req, err := http.NewRequest(r.Method, infocmdbUrl, bytes.NewReader(body))
		req, err := http.NewRequest(r.Method, infocmdbUrl, bytes.NewBufferString(values.Encode()))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}

		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		// req.PostForm = values

		httpClient := &http.Client{}
		resp, err := httpClient.Do(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}
		defer resp.Body.Close()

		w.WriteHeader(resp.StatusCode)
		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}
		fmt.Printf("Mocking:\n''''''\n%s\n''''''\n", string(respBody))
		w.Write(respBody)
	}))
}

func getUrl() string {
	testingURL := os.Getenv("WORKFLOW_TEST_URL")
	if mocking {
		ts := newMockServer()
		testingURL = ts.URL
	}

	if testingURL == "" {
		log.Fatal("WORKFLOW_TEST_URL must be provided or mocking enabled(WORKFLOW_TEST_MOCKING=true)")
	}

	log.Debugf("Testing-URL: %s", testingURL)
	return testingURL
}

func TestInfoCMDB_GetListOfCiIdsOfCiType(t *testing.T) {
	type fields struct {
		v1 *v1.Cmdb
		v2 *v2.Cmdb
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
				&v1.Cmdb{Config: v1.Config{}},
				&v2.Cmdb{Config: v2.Config{
					Url:      infocmdbUrl,
					Username: "fail",
					Password: "fail",
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
				&v1.Cmdb{Config: v1.Config{}},
				&v2.Cmdb{Config: v2.Config{
					Url:      infocmdbUrl,
					Username: "admin",
					Password: "admin",
					BasePath: "/app/",
				}},
			},
			args{ciTypeID: 1},
			CiIds{1, 2},
			false,
		},
		{
			"v2 List Ci's of Type '-1' (error)",
			fields{
				&v1.Cmdb{Config: v1.Config{}},
				&v2.Cmdb{Config: v2.Config{
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
			i := &Client{
				v1: tt.fields.v1,
				v2: tt.fields.v2,
			}
			gotR, err := i.GetListOfCiIdsOfCiType(tt.args.ciTypeID)
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
	type fields struct {
		v1 *v1.Cmdb
		v2 *v2.Cmdb
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
				&v1.Cmdb{Config: v1.Config{}},
				&v2.Cmdb{Config: v2.Config{
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
			"v2 List CIs of Type 1 (demo)",
			fields{
				&v1.Cmdb{Config: v1.Config{}},
				&v2.Cmdb{Config: v2.Config{
					Url:      infocmdbUrl,
					Username: "admin",
					Password: "admin",
					BasePath: "/app/",
				}},
			},
			args{ciTypeID: 1},
			CiIds{1, 2},
			false,
		},
		{
			"v2 List Ci's of Type '-1' (error)",
			fields{
				&v1.Cmdb{Config: v1.Config{}},
				&v2.Cmdb{Config: v2.Config{
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
			i := &Client{
				v1: tt.fields.v1,
				v2: tt.fields.v2,
			}
			gotR, err := i.GetListOfCiIdsOfCiTypeV2(tt.args.ciTypeID)
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

func TestInfoCMDB_QueryWebservice(t *testing.T) {
	type fields struct {
		v2 *v2.Cmdb
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
				&v2.Cmdb{Config: v2.Config{
					Url:      infocmdbUrl,
					Username: "admin",
					Password: "admin",
					BasePath: "/app/",
				}},
			},
			args{ws: "int_getCi", params: map[string]string{"argv1": "1"}},
			`{"success":true,"message":"Query executed successfully","data":[{"ci_id":"1","ci_type_id":"1","ci_type":"demo","project":"springfield","project_id":"4"}]}`,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Client{
				v2: tt.fields.v2,
			}
			gotR, err := i.QueryWebservice(tt.args.ws, tt.args.params)
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

func TestInfoCMDB_UpdateCiAttribute(t *testing.T) {
	type fields struct {
		v2 *v2.Cmdb
	}
	type args struct {
		ci int
		ua []UpdateCiAttribute
	}

	cmdbConfigValid := v2.Cmdb{Config: v2.Config{
		Url:      infocmdbUrl,
		Username: "admin",
		Password: "admin",
		BasePath: "/app/",
	}}
	baseCiID := 14 // ci to base this tests on
	baseAttributeName := "emp_lastname"

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"v2 Delete CI Attribute - fail - requires ciattributeid",
			fields{&cmdbConfigValid},
			args{ci: baseCiID, ua: []UpdateCiAttribute{
				{Mode: v2.UPDATE_MODE_DELETE, Name: baseAttributeName},
			}},
			true,
		},
		{
			"v2 Update CI Attribute",
			fields{&cmdbConfigValid},
			args{ci: baseCiID, ua: []UpdateCiAttribute{
				{Mode: v2.UPDATE_MODE_SET, Name: "emp_firstname", Value: "22322"},
			}},
			false,
		},
		{
			"v2 Update CI Attribute - wrong attribute name",
			fields{&cmdbConfigValid},
			args{ci: baseCiID, ua: []UpdateCiAttribute{
				{Mode: v2.UPDATE_MODE_SET, Name: baseAttributeName + "_NOT_EXISTING", Value: "1"},
			}},
			true,
		},
		{
			"v2 Insert CI Attribute",
			fields{&cmdbConfigValid},
			args{ci: baseCiID, ua: []UpdateCiAttribute{
				{Mode: v2.UPDATE_MODE_INSERT, Name: baseAttributeName, Value: "New1"},
			}},
			false,
		},
		{
			"v2 Update CI Attribute - fail - multiple attributes with the same name",
			fields{&cmdbConfigValid},
			args{ci: baseCiID, ua: []UpdateCiAttribute{
				{Mode: v2.UPDATE_MODE_SET, Name: baseAttributeName, Value: "22322"},
			}},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Client{
				v2: tt.fields.v2,
			}
			err := i.UpdateCiAttribute(tt.args.ci, tt.args.ua)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateCiAttribute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
