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
	"bytes"
	v1 "github.com/infonova/infocmdb-sdk-go/infocmdb/v1/infocmdb"
	v2 "github.com/infonova/infocmdb-sdk-go/infocmdb/v2/infocmdb"
	log "github.com/sirupsen/logrus"
	easy "github.com/t-tomalak/logrus-easy-formatter"
	"os"
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

type logOutputSplitter struct{}

func (splitter *logOutputSplitter) Write(msg []byte) (n int, err error) {
	if bytes.HasPrefix(msg, []byte("[TRACE]")) ||
		bytes.HasPrefix(msg, []byte("[DEBUG]")) ||
		bytes.HasPrefix(msg, []byte("[INFO]")) {
		return os.Stdout.Write(msg)
	}
	return os.Stderr.Write(msg)
}

func init() {
	// Log to stdout and stderr depending on log level:
	// Any message on stderr is interpreted as workflow failure
	log.SetOutput(&logOutputSplitter{})
	// Time is omitted in the log message, because it is already shown in the workflow log in a separate column
	log.SetFormatter(&easy.Formatter{
		LogFormat: "[%lvl%] %msg%\n",
	})
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
