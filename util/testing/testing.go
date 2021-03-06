package testing

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"

	"gopkg.in/yaml.v2"

	log "github.com/sirupsen/logrus"

	"github.com/joho/godotenv"
)

type Testing struct {
	mocking       bool
	mockings      map[string]mockingResponse
	mockingServer *httptest.Server
	url           string
}

func New() *Testing {
	t := Testing{}
	if os.Getenv("WORKFLOW_TEST_MOCKING") == "true" {
		t.mocking = true
		log.Debug("Mocking enabled")
	} else if err := godotenv.Load("../../.env"); err != nil {
		log.Fatalf("failed to load env: %v", err)
	}

	err := godotenv.Load("../../.env")
	//
	if err != nil {
		if t.mocking {
			log.Infof("ignoring failure to load env due to enabled Mocking, error: %v", err)
		} else {
			log.Fatalf("failed to load env: %v", err)
		}
	}

	t.setupMocking()

	testingURL := os.Getenv("WORKFLOW_TEST_URL")
	if t.mocking {
		mockServer := t.newMockServer()
		t.mockingServer = mockServer
		testingURL = mockServer.URL
	}

	if testingURL == "" {
		log.Fatal("WORKFLOW_TEST_URL must be provided or Mocking enabled(WORKFLOW_TEST_MOCKING=true)")
	}

	log.Debugf("Testing-URL: %s", testingURL)
	t.url = testingURL

	return &t
}

func (t *Testing) setupMocking() {
	t.mockings = make(map[string]mockingResponse)

	t.AddMocking(Mocking{
		RequestString: `GET##/##`,
		ReturnString:  ``,
	})

	// apiV1

	t.AddMocking(Mocking{
		RequestString: `GET##/api/login/username/admin/password/admin/timeout/600/method/json##`,
		ReturnString:  `{"status":"OK","apikey":"mocked_mocked_mocked_mocked###"}`,
	})
	t.AddMocking(Mocking{
		RequestString: `POST##/api/adapter/query/int_getListOfCiIdsOfCiType/method/json##apikey=mocked_mocked_mocked_mocked%23%23%23&argv1=1`,
		ReturnString:  `{"status":"OK","data":[{"ciid":"1"},{"ciid":"2"}]}`,
	})

	t.AddMocking(Mocking{
		RequestString: `PUT##/apiV2/query/execute/int_getCiAttributeId##{"query":{"params":{"argv1":"428","argv2":"29"}}}`,
		ReturnString:  `{"status":"OK","data":[{}]}`,
	})

	t.AddMocking(Mocking{
		RequestString: `GET##/apiV2/ci/index?ciTypeId=-1##`,
		ReturnString:  `{"success":false,"message":"Internal Server Error","data":null}`,
		StatusCode:    http.StatusInternalServerError,
	})

	// apiV2
	t.AddMocking(Mocking{
		RequestString: `POST##/apiV2/auth/token##lifetime=600&password=fail&username=fail`,
		StatusCode:    http.StatusUnauthorized,
	})

	t.AddMocking(Mocking{
		RequestString: `POST##/apiV2/auth/token##lifetime=600&password=false&username=false`,
		StatusCode:    http.StatusBadRequest,
	})

	t.AddMocking(Mocking{
		RequestString: `POST##/apiV2/auth/token##lifetime=600&password=false&username=admin`,
		StatusCode:    http.StatusBadRequest,
	})

	t.AddMocking(Mocking{
		RequestString: `POST##/apiV2/auth/token##lifetime=600&password=fail&username=admin`,
		StatusCode:    http.StatusBadRequest,
	})

	t.AddMocking(Mocking{
		RequestString: `POST##/apiV2/auth/token##lifetime=600&password=admin&username=admin`,
		ReturnString:  `{"success":true,"message":"mocked","data":{"token":"mocked"}}`,
	})

	t.AddMocking(Mocking{
		RequestString: `PUT##/apiV2/query/execute/int_createCi##{"query":{"params":{"argv1":"476","argv2":"","argv3":"0"}}}`,
		ReturnString:  `{"success":true,"message":"Query executed successfully","data":[{"id":"617827","ci_type_id":"476","icon":"","history_id":"59529024","valid_from":"2020-01-13 15:14:05","created_at":"2020-01-13 15:14:05","updated_at":"2020-01-13 15:14:05"}]}`,
	})

	t.AddMocking(Mocking{
		RequestString: `PUT##/apiV2/query/execute/int_addCiProjectMapping##{"query":{"params":{"argv1":"617830","argv2":"33","argv3":"59529029"}}}`,
		ReturnString:  `{"success":true,"message":"Query executed successfully","data":null}`,
	})

	t.AddMocking(Mocking{
		RequestString: `PUT##/apiV2/query/execute/int_getListOfCiIdsOfCiType##{"query":{"params":{"argv1":"1"}}}`,
		ReturnString:  `{"status":"OK","data":[{"ciid":"1"},{"ciid":"2"}]}`,
	})

	t.AddMocking(Mocking{
		RequestString: `PUT##/apiV2/query/execute/int_getCi##{"query":{"params":{"argv1":"1"}}}`,
		ReturnString:  `{"success":true,"message":"Query executed successfully","data":[{"ci_id":"1","ci_type_id":"1","ci_type":"demo","project":"springfield","project_id":"4"}]}`,
	})

	t.AddMocking(Mocking{
		RequestString: `PUT##/apiV2/ci/14##{"ci":{"attributes":[{"mode":"set","name":"emp_firstname","value":"22322","ciAttributeId":0,"uploadId":""}]}}`,
		ReturnString:  `{"success":true,"message":"Query executed successfully","data":[]}`,
	})

	t.AddMocking(Mocking{
		RequestString: `PUT##/apiV2/ci/14##{"ci":{"attributes":[{"mode":"insert","name":"emp_lastname","value":"New1","ciAttributeId":0,"uploadId":""}]}}`,
		ReturnString:  `{"success":true,"message":"Query executed successfully","data":[]}`,
	})

	t.AddMocking(Mocking{
		RequestString: `PUT##/apiV2/ci/14##{"ci":{"attributes":[{"mode":"delete","name":"emp_lastname","value":"","ciAttributeId":0,"uploadId":""}]}}`,
		StatusCode:    http.StatusBadRequest,
	})

	t.AddMocking(Mocking{
		RequestString: `PUT##/apiV2/ci/14##{"ci":{"attributes":[{"mode":"set","name":"emp_lastname_NOT_EXISTING","value":"1","ciAttributeId":0,"uploadId":""}]}}`,
		StatusCode:    http.StatusBadRequest,
	})

	t.AddMocking(Mocking{
		RequestString: `PUT##/apiV2/ci/14##{"ci":{"attributes":[{"mode":"set","name":"emp_lastname","value":"22322","ciAttributeId":0,"uploadId":""}]}}`,
		StatusCode:    http.StatusBadRequest,
	})

	t.AddMocking(Mocking{
		RequestString: `PUT##/apiV2/query/execute/int_getAttributeDefaultOptionId##{"query":{"params":{"argv1":"438","argv2":"IN PROGRESS"}}}`,
		ReturnString:  `{"status":"OK","data":[{"id":"1329"}]}`,
	})
}

func (t *Testing) GetUrl() string {
	return t.url
}

type Mocking struct {
	RequestString string
	ReturnString  string
	ContentType   string
	StatusCode    int
}

type mockingResponse struct {
	ReturnString string
	ContentType  string
	StatusCode   int
}

func (t *Testing) AddMocking(m Mocking) *Testing {
	t.mockings[m.RequestString] = mockingResponse{
		ReturnString: m.ReturnString,
		ContentType:  m.ContentType,
		StatusCode:   m.StatusCode,
	}

	return t
}

func (t *Testing) newMockServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		mockString := fmt.Sprintf("%s##%s##%s", r.Method, r.URL.String(), string(body))

		if m, ok := t.mockings[mockString]; ok {
			if m.ContentType == "" {
				m.ContentType = "application/json;charset=UTF-8"
			}
			if m.StatusCode == 0 {
				m.StatusCode = http.StatusOK
			}

			w.Header()["Content-Type"] = []string{m.ContentType}
			w.WriteHeader(m.StatusCode)
			if _, err := w.Write([]byte(m.ReturnString)); err != nil {
				log.Error(err)
			}
			return
		} else {
			log.Fatalf("Didn't Mock:\n''''''\n%s\n''''''\n", mockString)
		}
	}))
}

func (t *Testing) SetValidConfig(config interface{}) {
	configBytes := []byte(fmt.Sprintf(`version: 1.0
apiUrl: %v
apiUser: admin
apiPassword: admin
`, t.GetUrl()))
	err := yaml.Unmarshal(configBytes, config)
	if err != nil {
		log.Fatalf("failed to build valid config: %v", err)
	}
}
