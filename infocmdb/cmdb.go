package infocmdb

import (
	v1 "github.com/infonova/infocmdb-lib-go/core/v1/cmdb"
	v2 "github.com/infonova/infocmdb-lib-go/core/v2/cmdb"
	log "github.com/sirupsen/logrus"
	"os"
)

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

func init() {
	log.SetLevel(log.InfoLevel)
	if os.Getenv("WORKFLOW_DEBUGGING") == "true" {
		log.SetLevel(log.DebugLevel)
	}
}

// NewClient returns a new cmdb client
func NewClient() *Client {
	return new(Client)
}

// LoadConfig from file in yaml format
func (c *Client) LoadConfig(f string) (err error) {
	c.v1, err = v1.New(f)
	if err != nil {
		return
	}

	c.v2 = v2.New()
	if err = c.v2.LoadConfigFile(f); err != nil {
		return
	}

	return
}
