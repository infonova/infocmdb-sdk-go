package infocmdblibrary

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
)

type testing struct{}

var infocmdbCredentials = Credentials{Username: "admin", Password: "admin"}

func (t *testing) MockServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		url := "http://infocmdb.local" + r.URL.String()
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
			w.Write([]byte(`{"status":"OK","data":[{"ci_id":"1","ci_attribute_id":"1","attribute_id":"1","attribute_name":"general_unique_input","attribute_description":"Unique Input","attribute_type":"input","value":"demo_1","modified_at":"2018-11-16 08:13:21"},{"ci_id":"1","ci_attribute_id":"2","attribute_id":"3","attribute_name":"general_regular_input","attribute_description":"Regular Input","attribute_type":"input","value":"Regular Single Line Text Input","modified_at":"2018-11-16 08:13:21"},{"ci_id":"1","ci_attribute_id":"3","attribute_id":"4","attribute_name":"general_numeric_input","attribute_description":"Numeric Input","attribute_type":"input","value":"42","modified_at":"2018-11-16 08:13:21"},{"ci_id":"1","ci_attribute_id":"4","attribute_id":"5","attribute_name":"general_textarea","attribute_description":"Textarea","attribute_type":"textarea","value":"Multiline Text Input\n\n--> MORE TEXT","modified_at":"2018-11-16 08:13:21"},{"ci_id":"1","ci_attribute_id":"5","attribute_id":"6","attribute_name":"general_textedit","attribute_description":"Editor Area","attribute_type":"textEdit","value":"Multiline Text Input<br \/><br \/>WITH <i>STYLING<\/i> <b>OPTIONS<\/b>!","modified_at":"2018-11-16 08:13:21"},{"ci_id":"1","ci_attribute_id":"6","attribute_id":"7","attribute_name":"general_dropdown_static","attribute_description":"Dropdown (static)","attribute_type":"select","value":"Option 1","modified_at":"2018-11-16 08:13:21"},{"ci_id":"1","ci_attribute_id":"7","attribute_id":"8","attribute_name":"general_checkbox","attribute_description":"Checkbox","attribute_type":"checkbox","value":"Check 3","modified_at":"2018-11-16 08:13:21"},{"ci_id":"1","ci_attribute_id":"8","attribute_id":"9","attribute_name":"general_radio","attribute_description":"Radio","attribute_type":"radio","value":"Radio 2","modified_at":"2018-11-16 08:13:21"},{"ci_id":"1","ci_attribute_id":"9","attribute_id":"10","attribute_name":"general_date","attribute_description":"Date","attribute_type":"date","value":"2019-02-18","modified_at":"2018-11-16 08:13:21"},{"ci_id":"1","ci_attribute_id":"10","attribute_id":"11","attribute_name":"general_datetime","attribute_description":"Datetime","attribute_type":"dateTime","value":"2013-11-01 12:24:10","modified_at":"2018-11-16 08:13:21"},{"ci_id":"1","ci_attribute_id":"11","attribute_id":"12","attribute_name":"general_currency","attribute_description":"Currency","attribute_type":"zahlungsmittel","value":"29,99","modified_at":"2018-11-16 08:13:21"},{"ci_id":"1","ci_attribute_id":"12","attribute_id":"13","attribute_name":"general_password","attribute_description":"Password","attribute_type":"password","value":"secret demo password","modified_at":"2018-11-16 08:13:21"},{"ci_id":"1","ci_attribute_id":"13","attribute_id":"14","attribute_name":"general_hyperlink","attribute_description":"Hyperlink","attribute_type":"link","value":"#todo","modified_at":"2018-11-16 08:13:21"},{"ci_id":"1","ci_attribute_id":"14","attribute_id":"23","attribute_name":"general_dropdown_sql_filled_select","attribute_description":"Dropdown (SQL filled) - Regular Input","attribute_type":"selectQuery","value":"#todo","modified_at":"2018-11-16 08:13:21"},{"ci_id":"1","ci_attribute_id":"15","attribute_id":"25","attribute_name":"general_dropdown_sql_filled_multiselect","attribute_description":"Dropdown (SQL filled) - Multiselect","attribute_type":"selectQuery","value":"#todo","modified_at":"2018-11-16 08:13:21"},{"ci_id":"1","ci_attribute_id":"16","attribute_id":"26","attribute_name":"general_dropdown_sql_filled_multiselect_counter","attribute_description":"Dropdown (SQL filled) - Multiselect Counter","attribute_type":"selectQuery","value":"#todo","modified_at":"2018-11-16 08:13:21"}]}`))
			return
		case "GET##/api/adapter/apikey/4afbf95c1d072664e35cd61339e152/query/int_getListOfCiIdsOfCiType/method/json?apikey=4afbf95c1d072664e35cd61339e152&argv1=1##":
			w.WriteHeader(200)
			w.Write([]byte(`{"status":"OK","data":[{"ciid":"1"},{"ciid":"2"}]}`))
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
		proxyReq, err := http.NewRequest(r.Method, url, bytes.NewReader(body))
		proxyReq.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		httpClient := http.Client{}
		resp, err := httpClient.Do(proxyReq)
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
		fmt.Println(string(respBody))
		w.Write(respBody)
	}))
}

func (t *testing) getUrl() string {
	if os.Getenv("testing_mocking") == "true" {
		ts := t.MockServer()
		return ts.URL
	}

	return os.Getenv("testing_url")
}

func ExampleWebservice_Webservice() {
	t := testing{}
	infocmdbURL := t.getUrl()

	i, err := NewCMDB(infocmdbURL, infocmdbCredentials)
	if err != nil {

		fmt.Printf("Error: %v\n", err)
		return
	}

	params := url.Values{}
	params.Add("argv1", "1")
	ret, err := i.WS.Webservice("int_getListOfCiIdsOfCiType", params)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Return: %v\n", ret)

	// Output:
	// Return: {"status":"OK","data":[{"ciid":"1"},{"ciid":"2"}]}
}

func ExampleCmdbWebClient_Login() {
	t := testing{}
	infocmdbURL := t.getUrl()

	i, err := NewCMDB(infocmdbURL, infocmdbCredentials)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Login ok, ApiKey(len): %d\n", len(i.WS.client.apikey))

	// Output:
	// Login ok, ApiKey(len): 30
}

func ExampleCmdbWebClient_Get() {
	t := testing{}
	infocmdbURL := t.getUrl()

	i, err := NewCMDB(infocmdbURL, infocmdbCredentials)
	if err != nil {

		fmt.Printf("Error: %v\n", err)
		return
	}

	params := url.Values{
		"argv1": {"1"},
	}

	ret, err := i.WC.Get("query", "int_getListOfCiIdsOfCiType", params)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Println("Get: ", ret)

	// Output:
	// Get:  {"status":"OK","data":[{"ciid":"1"},{"ciid":"2"}]}
}

func ExampleCmdbWebClient_Post() {
	t := testing{}
	infocmdbURL := t.getUrl()

	i, err := NewCMDB(infocmdbURL, infocmdbCredentials)
	if err != nil {

		fmt.Printf("Error: %v\n", err)
		return
	}

	params := url.Values{
		"argv1": {"1"},
	}

	ret, err := i.WC.Post("query", "int_getListOfCiIdsOfCiType", params)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Println("Post: ", ret)

	// Output:
	// Post:  {"status":"OK","data":[{"ciid":"1"},{"ciid":"2"}]}
}

func ExampleInfoCmdbGoLib_GetCi() {
	t := testing{}
	infocmdbURL := t.getUrl()

	i, err := NewCMDB(infocmdbURL, infocmdbCredentials)
	if err != nil {

		fmt.Printf("Error: %v\n", err)
		return
	}

	rCi, err := i.GetCi(1)
	if err != nil {
		fmt.Printf("WS Error: %v", err)
		return
	}
	fmt.Println(rCi)

	// Output:
	// {OK [{1 1 demo springfield 4}]}
}

func ExampleInfoCmdbGoLib_GetListOfCiIdsOfCiType() {
	t := testing{}
	infocmdbURL := t.getUrl()

	i, err := NewCMDB(infocmdbURL, infocmdbCredentials)
	if err != nil {

		fmt.Printf("Error: %v\n", err)
		return
	}

	rCi, err := i.GetListOfCiIdsOfCiType(1)
	if err != nil {
		fmt.Printf("WS Error: %v\n", err)
		return
	}
	fmt.Println(rCi)

	// Output:
	// {OK [{1} {2}]}
}

func ExampleInfoCmdbGoLib_Login() {
	t := testing{}
	infocmdbURL := t.getUrl()

	i, err := NewCMDB(infocmdbURL, infocmdbCredentials)
	if err != nil {

		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Println("Apikey length: ", len(i.WC.apikey))

	// Output:
	// Apikey length:  30
}

func ExampleInfoCmdbGoLib_GetCiAttributes() {
	t := testing{}
	infocmdbURL := t.getUrl()

	i, err := NewCMDB(infocmdbURL, infocmdbCredentials)
	if err != nil {

		fmt.Printf("Error: %v\n", err)
		return
	}

	rCi, err := i.GetCiAttributes(1)
	if err != nil {
		fmt.Printf("ws %v\n", err)
		return
	}
	fmt.Println(rCi)

	// Output:
	// {OK [{1 1 1 general_unique_input Unique Input input demo_1 2018-11-16 08:13:21} {1 2 3 general_regular_input Regular Input input Regular Single Line Text Input 2018-11-16 08:13:21} {1 3 4 general_numeric_input Numeric Input input 42 2018-11-16 08:13:21} {1 4 5 general_textarea Textarea textarea Multiline Text Input
	//
	//--> MORE TEXT 2018-11-16 08:13:21} {1 5 6 general_textedit Editor Area textEdit Multiline Text Input<br /><br />WITH <i>STYLING</i> <b>OPTIONS</b>! 2018-11-16 08:13:21} {1 6 7 general_dropdown_static Dropdown (static) select Option 1 2018-11-16 08:13:21} {1 7 8 general_checkbox Checkbox checkbox Check 3 2018-11-16 08:13:21} {1 8 9 general_radio Radio radio Radio 2 2018-11-16 08:13:21} {1 9 10 general_date Date date 2019-02-18 2018-11-16 08:13:21} {1 10 11 general_datetime Datetime dateTime 2013-11-01 12:24:10 2018-11-16 08:13:21} {1 11 12 general_currency Currency zahlungsmittel 29,99 2018-11-16 08:13:21} {1 12 13 general_password Password password secret demo password 2018-11-16 08:13:21} {1 13 14 general_hyperlink Hyperlink link #todo 2018-11-16 08:13:21} {1 14 23 general_dropdown_sql_filled_select Dropdown (SQL filled) - Regular Input selectQuery #todo 2018-11-16 08:13:21} {1 15 25 general_dropdown_sql_filled_multiselect Dropdown (SQL filled) - Multiselect selectQuery #todo 2018-11-16 08:13:21} {1 16 26 general_dropdown_sql_filled_multiselect_counter Dropdown (SQL filled) - Multiselect Counter selectQuery #todo 2018-11-16 08:13:21}]}
}
