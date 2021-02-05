package infocmdb

import (
	"fmt"
	"net/http"
	"net/url"

	log "github.com/sirupsen/logrus"

	utilTesting "github.com/infonova/infocmdb-sdk-go/util/testing"
)

func ExampleCmdb_Webservice_int_getListOfCiIdsOfCiType() {
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

func ExampleCmdb_LoadConfigFile_fail() {
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

func ExampleCmdb_Login() {
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
func ExampleCmdb_Login_with_ApiKey() {
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

func ExampleCmdb_CallWebservice_post_query() {
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
