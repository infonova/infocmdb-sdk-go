package infocmdblibrary

import (
	"fmt"
	"net/url"
)

var infocmdbURL = "http://infocmdb.local"
var infocmdbCredentials = Credentials{Username: "admin", Password: "admin"}

func ExampleWebservice_Webservice() {

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
