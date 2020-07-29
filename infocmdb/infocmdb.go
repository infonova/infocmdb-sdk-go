package infocmdb

// This library is used for the communication with the infoCMDB
//
// The API provided by the infoCMDB has different versions and therefore this library is split up v1/v2.
//
// v1 - custom HTTP API
//
// Is the legacy version which is based on configured sql queries provided via a custom http api
//
// v2 - Restful API
//
// This is the first re-engineering or the api to access core models and services to have native access.
// This api properly handles all permission checks and access to native functions.

import (
	v1 "github.com/infonova/infocmdb-sdk-go/infocmdb/v1/infocmdb"
	v2 "github.com/infonova/infocmdb-sdk-go/infocmdb/v2/infocmdb"
)

// Client configuration values.
// Usually taken from data/configs/workflows/infocmdb.yml
type Config struct {
	ApiUrl       string `yaml:"apiUrl"`
	ApiUser      string `yaml:"apiUser"`
	ApiPassword  string `yaml:"apiPassword"`
	ApiKey       string
	CmdbBasePath string `yaml:"CmdbBasePath"`
}

// Client combines connectivity methods for version 1 and 2 of the cmdb
type Client struct {
	v1 *v1.Cmdb
	v2 *v2.Cmdb
}

// NewClient returns a new cmdb client
func NewClient() (c *Client) {
	c = &Client{
		v1: v1.New(),
		v2: v2.New(),
	}
	return
}

// LoadConfig from file in yaml format
func (c *Client) LoadConfig(path string) (err error) {
	err = c.v1.LoadConfigFile(path)
	if err != nil {
		return
	}

	err = c.v2.LoadConfigFile(path)
	if err != nil {
		return
	}

	return
}

// Create a client for use in unit tests (possibly of other libraries that wrap infocmdb-sdk-go).
// A default configuration with username and password "admin" is used.
func NewTestClient(url string) (c *Client) {
	c = NewClient()
	c.v1.LoadConfig(v1.Config{
		ApiUrl:      url,
		ApiUser:     "admin",
		ApiPassword: "admin",
	})
	c.v2.LoadConfig(v2.Config{
		Url:      url,
		Username: "admin",
		Password: "admin",
	})
	return
}
