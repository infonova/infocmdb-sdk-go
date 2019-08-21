package cmdb

import (
	"bytes"
	"fmt"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
)

type testing struct{}

var (
	ErrTestingInfocmdbUrlMissing = "WORKFLOW_TEST_URL must be provided or mocking enabled(INFOCMDB_WORKFLOW_TEST_MOCKING=true)"
)

var (
	mocking            = false
	infoCMDBConfig = []byte(`apiUrl: http://nginx/
apiUser: admin
apiPassword: admin
CmdbBasePath: /app/`)
	infoCMDBConfigFile = "test/test.yml"
)

func init() {
	if err := godotenv.Load("../../.env"); err != nil {
		log.Fatalf("failed to load env: %v", err)
	}

	if os.Getenv("WORKFLOW_TEST_MOCKING") == "true" {
		mocking = true
		log.Debug("Mocking enabled")
	}
}

func (t *testing) MockServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		infocmdbUrl := os.Getenv("WORKFLOW_TEST_URL") + r.URL.String()
		mockString := fmt.Sprintf("%s##%s##%s", r.Method, r.URL.String(), string(body))
		switch mockString {
		case "GET##/api/login/username/admin/password/admin/timeout/21600/method/json##":
			w.WriteHeader(200)
			w.Write([]byte(`{"status":"OK","apikey":"4afbf95c1d072664e35cd61339e152"}`))
			return
		case "POST##/api/adapter/query/int_getListOfCiIdsOfCiType/method/json##apikey=4afbf95c1d072664e35cd61339e152&argv1=1":
			w.WriteHeader(200)
			w.Write([]byte(`{"status":"OK","data":[{"ciid":"1"},{"ciid":"2"}]}`))
			return
		case "POST##/api/adapter/query/int_getCi/method/json##apikey=4afbf95c1d072664e35cd61339e152&argv1=1":
			w.WriteHeader(200)
			w.Write([]byte(`{"status":"OK","data":[{"ci_id":"1","ci_type_id":"1","ci_type":"demo","project":"springfield","project_id":"4"}]}`))
			return
		case "POST##/api/adapter/query/int_getCiAttributes/method/json##apikey=4afbf95c1d072664e35cd61339e152&argv1=1":
			w.WriteHeader(200)
			w.Write([]byte(`{"status":"OK","data":[{"ci_id":"1","ci_attribute_id":"1","attribute_id":"1","attribute_name":"general_unique_input","attribute_description":"Unique Input","attribute_type":"input","value":"demo_1","modified_at":"2019-01-30 12:16:54"},{"ci_id":"1","ci_attribute_id":"2","attribute_id":"3","attribute_name":"general_regular_input","attribute_description":"Regular Input","attribute_type":"input","value":"Regular Single Line Text Input","modified_at":"2019-01-3012:16:54"},{"ci_id":"1","ci_attribute_id":"3","attribute_id":"4","attribute_name":"general_numeric_input","attribute_description":"Numeric Input","attribute_type":"input","value":"42","modified_at":"2019-01-30 12:16:54"},{"ci_id":"1","ci_attribute_id":"4","attribute_id":"5","attribute_name":"general_textarea","attribute_description":"Textarea","attribute_type":"textarea","value":"Multiline Text Input\n\n--> MORE TEXT","modified_at":"2019-01-30 12:16:54"},{"ci_id":"1","ci_attribute_id":"5","attribute_id":"6","attribute_name":"general_textedit","attribute_description":"Editor Area","attribute_type":"textEdit","value":"Multiline Text Input<br \/><br \/>WITH <i>STYLING<\/i> <b>OPTIONS<\/b>!","modified_at":"2019-01-30 12:16:54"},{"ci_id":"1","ci_attribute_id":"6","attribute_id":"7","attribute_name":"general_dropdown_static","attribute_description":"Dropdown (static)","attribute_type":"select","value":"Option 1","modified_at":"2019-01-30 12:16:54"},{"ci_id":"1","ci_attribute_id":"7","attribute_id":"8","attribute_name":"general_checkbox","attribute_description":"Checkbox","attribute_type":"checkbox","value":"Check 3","modified_at":"2019-01-30 12:16:54"},{"ci_id":"1","ci_attribute_id":"8","attribute_id":"9","attribute_name":"general_radio","attribute_description":"Radio","attribute_type":"radio","value":"Radio 2","modified_at":"2019-01-30 12:16:54"},{"ci_id":"1","ci_attribute_id":"9","attribute_id":"10","attribute_name":"general_date","attribute_description":"Date","attribute_type":"date","value":"2019-02-18","modified_at":"2019-01-30 12:16:54"},{"ci_id":"1","ci_attribute_id":"10","attribute_id":"11","attribute_name":"general_datetime","attribute_description":"Datetime","attribute_type":"dateTime","value":"2013-11-01 12:24:10","modified_at":"2019-01-30 12:16:54"},{"ci_id":"1","ci_attribute_id":"11","attribute_id":"12","attribute_name":"general_currency","attribute_description":"Currency","attribute_type":"zahlungsmittel","value":"29,99","modified_at":"2019-01-30 12:16:54"},{"ci_id":"1","ci_attribute_id":"12","attribute_id":"13","attribute_name":"general_password","attribute_description":"Password","attribute_type":"password","value":"secret demo password","modified_at":"2019-01-30 12:16:54"},{"ci_id":"1","ci_attribute_id":"13","attribute_id":"14","attribute_name":"general_hyperlink","attribute_description":"Hyperlink","attribute_type":"link","value":"#todo","modified_at":"2019-01-30 12:16:54"},{"ci_id":"1","ci_attribute_id":"14","attribute_id":"16","attribute_name":"general_regular_executable","attribute_description":"Executable Script","attribute_type":"executeable","value":null,"modified_at":"2019-01-30 12:16:54"},{"ci_id":"1","ci_attribute_id":"15","attribute_id":"17","attribute_name":"general_event_executable","attribute_description":"Event Executable","attribute_type":"executeable","value":null,"modified_at":"2019-01-30 12:16:54"},{"ci_id":"1","ci_attribute_id":"16","attribute_id":"23","attribute_name":"general_dropdown_sql_filled_select","attribute_description":"Dropdown (SQL filled) - Regular Input","attribute_type":"selectQuery","value":"#todo","modified_at":"2019-01-30 12:16:54"},{"ci_id":"1","ci_attribute_id":"17","attribute_id":"25","attribute_name":"general_dropdown_sql_filled_multiselect","attribute_description":"Dropdown (SQL filled) - Multiselect","attribute_type":"selectQuery","value":"#todo","modified_at":"2019-01-3012:16:54"},{"ci_id":"1","ci_attribute_id":"18","attribute_id":"26","attribute_name":"general_dropdown_sql_filled_multiselect_counter","attribute_description":"Dropdown (SQL filled) - Multiselect Counter","attribute_type":"selectQuery","value":"#todo","modified_at":"2019-01-30 12:16:54"}]}`))
			return
		case "GET##/api/adapter/apikey/4afbf95c1d072664e35cd61339e152/query/int_getListOfCiIdsOfCiType/method/json?apikey=4afbf95c1d072664e35cd61339e152&argv1=1##":
			w.WriteHeader(200)
			w.Write([]byte(`{"status":"OK","data":[{"ciid":"1"},{"ciid":"2"}]}`))
			return
		case "POST##/api/adapter/query/int_getCi/method/json##apikey=4afbf95c1d072664e35cd61339e152&argv1=-1":
			w.WriteHeader(200)
			w.Write([]byte(`{"status":"OK","data":[]}`))
			return
		// case "":
		// 	w.WriteHeader(200)
		// 	w.Write([]byte(``))
		// 	return
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

func (t *testing) getUrl() string {
	testingURL := os.Getenv("WORKFLOW_TEST_URL")
	if mocking {
		ts := t.MockServer()
		testingURL = ts.URL
	}

	if testingURL == "" {
		log.Fatal(ErrTestingInfocmdbUrlMissing)
	}

	log.Debugf("Testing-URL: %s", testingURL)
	return testingURL
}

func ExampleWebservice_Webservice() {
	t := testing{}

	i, err := NewCMDB(infoCMDBConfigFile)
	i.Config.ApiUrl = t.getUrl()
	if err != nil {
		log.Error(ErrFailedToCreateInfoCMDB)
		return
	}

	err = i.Login()
	if err != nil {
		log.Error(err)
		return
	}

	params := url.Values{}
	params.Add("argv1", "1")
	ret, err := i.Webservice("int_getListOfCiIdsOfCiType", params)
	if err != nil {
		log.Error(err)
		return
	}
	fmt.Printf("Return: %v\n", ret)

	// Output:
	// Return: {"status":"OK","data":[{"ciid":"1"},{"ciid":"2"}]}
}

func ExampleInfoCmdbGoLib_LoadConfigAbsolutePath() {
	i := InfoCMDB{}
	err := i.LoadConfig(infoCMDBConfig)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	fmt.Printf("Config: %v\n", i.Config)
	fmt.Printf("BasePath: %s\n", i.Config.CmdbBasePath)
	// Output:
	// Config: {http://nginx/ admin admin  /app/}
	// BasePath: /app/
}

func ExampleInfoCmdbGoLib_LoadConfig() {
	i := InfoCMDB{}
	err := i.LoadConfig(infoCMDBConfig)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	fmt.Printf("Config: %v\n", i.Config)
	fmt.Printf("BasePath: %s\n", i.Config.CmdbBasePath)
	// Output:
	// Config: {http://nginx/ admin admin  /app/}
	// BasePath: /app/
}

func ExampleInfoCmdbGoLib_LoadConfig_Fail() {
	i := InfoCMDB{}
	err := i.LoadConfigFile("/test_missing.yml")
	if err != nil {
		fmt.Println("loading failed")
		return
	}
	fmt.Printf("Config: %v\n", i.Config)
	fmt.Printf("BasePath: %s\n", i.Config.CmdbBasePath)
	// Output:
	// loading failed

}

func ExampleCmdbWebClient_Login() {
	t := testing{}

	i, err := NewCMDB(infoCMDBConfigFile)
	i.Config.ApiUrl = t.getUrl()
	if i == nil {
		log.Error(ErrFailedToCreateInfoCMDB)
		return
	}

	i.Config.ApiKey = ""

	err = i.Login()
	if err != nil {
		log.Error(err)
		return
	}

	fmt.Printf("Login ok, ApiKey(len): %d\n", len(i.Config.ApiKey))

	// Output:
	// Login ok, ApiKey(len): 30
}
func ExampleCmdbWebClient_LoginWithApiKey() {
	t := testing{}

	ilogin, err := NewCMDB(infoCMDBConfigFile)
	if err != nil {
		log.Error(err)
		return
	}

	ilogin.Config.ApiUrl = t.getUrl()
	err = ilogin.Login()
	if err != nil {
		log.Error(err)
		return
	}

	log.Debugf("Got API Key: %s", ilogin.Config.ApiKey)
	i, err := NewCMDB(infoCMDBConfigFile)
	if err != nil {
		log.Error(err)
		return
	}

	i.Config.ApiKey = ilogin.Config.ApiKey
	i.Config.ApiUrl = t.getUrl()

	fmt.Printf("Login ok, ApiKey(len): %d\n", len(i.Config.ApiKey))

	// Output:
	// Login ok, ApiKey(len): 30
}

func ExampleCmdbWebClient_Post() {
	t := testing{}

	i, err := NewCMDB(infoCMDBConfigFile)
	if err != nil {
		log.Error(err)
		return
	}

	i.Config.ApiUrl = t.getUrl()
	if i == nil {
		log.Error(ErrFailedToCreateInfoCMDB)
		return
	}

	params := url.Values{
		"argv1": {"1"},
	}

	ret := ""
	err = i.CallWebservice(http.MethodPost, "query", "int_getListOfCiIdsOfCiType", params, &ret)
	if err != nil {
		log.Error(err)
		return
	}
	fmt.Println("Post: ", ret)

	// Output:
	// Post:  {"status":"OK","data":[{"ciid":"1"},{"ciid":"2"}]}
}

