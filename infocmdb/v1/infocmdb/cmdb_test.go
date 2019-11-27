package infocmdb

import (
	"fmt"
	"net/http"
	"net/url"

	log "github.com/sirupsen/logrus"

	utilTesting "github.com/infonova/infocmdb-sdk-go/util/testing"
)

var (
	infoCMDBConfig = []byte(`apiUrl: http://nginx/
apiUser: admin
apiPassword: admin
CmdbBasePath: /app/`)
	infoCMDBConfigFile = "test/test.yml"
)

func ExampleWebservice_Webservice() {

	i, err := New(infoCMDBConfigFile)
	if err != nil {
		log.Error(ErrFailedToCreateInfoCMDB)
		return
	}

	i.Config.ApiUrl = utilTesting.New().GetUrl()

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
	i := Cmdb{}
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
	i := Cmdb{}
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

	i, err := New(infoCMDBConfigFile)
	if i == nil {
		log.Error(ErrFailedToCreateInfoCMDB)
		return
	}

	i.Config.ApiUrl = utilTesting.New().GetUrl()
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

	ilogin, err := New(infoCMDBConfigFile)
	if err != nil {
		log.Error(err)
		return
	}

	ilogin.Config.ApiUrl = utilTesting.New().GetUrl()
	err = ilogin.Login()
	if err != nil {
		log.Error(err)
		return
	}

	log.Debugf("Got API Key: %s", ilogin.Config.ApiKey)
	i, err := New(infoCMDBConfigFile)
	if err != nil {
		log.Error(err)
		return
	}

	i.Config.ApiKey = ilogin.Config.ApiKey
	i.Config.ApiUrl = utilTesting.New().GetUrl()

	fmt.Printf("Login ok, ApiKey(len): %d\n", len(i.Config.ApiKey))

	// Output:
	// Login ok, ApiKey(len): 30
}

func ExampleCmdbWebClient_Post() {

	i, err := New(infoCMDBConfigFile)
	if err != nil {
		log.Error(err)
		return
	}

	i.Config.ApiUrl = utilTesting.New().GetUrl()
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
