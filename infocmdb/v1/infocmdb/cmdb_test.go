package infocmdb

import (
	"fmt"
	"net/http"
	"net/url"

	log "github.com/sirupsen/logrus"

	utilTesting "github.com/infonova/infocmdb-sdk-go/util/testing"
)

func ExampleWebservice_Webservice() {

	cmdb := Cmdb{}
	utilTesting.New().SetValidConfig(&cmdb.Config)

	err := cmdb.Login()
	if err != nil {
		log.Error(err)
		return
	}

	params := url.Values{}
	params.Add("argv1", "1")
	ret, err := cmdb.Webservice("int_getListOfCiIdsOfCiType", params)
	if err != nil {
		log.Error(err)
		return
	}
	fmt.Printf("Return: %v\n", ret)

	// Output:
	// Return: {"status":"OK","data":[{"ciid":"1"},{"ciid":"2"}]}
}

func ExampleInfoCmdbGoLib_LoadConfig_Fail() {
	i := Cmdb{}
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

	cmdb := Cmdb{}
	utilTesting.New().SetValidConfig(&cmdb.Config)

	err := cmdb.Login()
	if err != nil {
		log.Error(err)
		return
	}

	fmt.Printf("Login ok, ApiKey(len): %d\n", len(cmdb.Config.ApiKey))

	// Output:
	// Login ok, ApiKey(len): 30
}
func ExampleCmdbWebClient_LoginWithApiKey() {

	cmdb := Cmdb{}
	utilTesting.New().SetValidConfig(&cmdb.Config)

	err := cmdb.Login()
	if err != nil {
		log.Error(err)
		return
	}

	fmt.Printf("Login ok, ApiKey(len): %d\n", len(cmdb.Config.ApiKey))

	// Output:
	// Login ok, ApiKey(len): 30
}

func ExampleCmdbWebClient_Post() {

	cmdb := Cmdb{}
	utilTesting.New().SetValidConfig(&cmdb.Config)

	params := url.Values{
		"argv1": {"1"},
	}

	ret := ""
	err := cmdb.CallWebservice(http.MethodPost, "query", "int_getListOfCiIdsOfCiType", params, &ret)
	if err != nil {
		log.Error(err)
		return
	}
	fmt.Println("Post: ", ret)

	// Output:
	// Post:  {"status":"OK","data":[{"ciid":"1"},{"ciid":"2"}]}
}
