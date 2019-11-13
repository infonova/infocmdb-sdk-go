package cmdb

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/infonova/infocmdb-sdk-go/infocmdb/v2/cmdb/client"
	"github.com/patrickmn/go-cache"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Url      string `yaml:"apiUrl"`
	Username string `yaml:"apiUser"`
	Password string `yaml:"apiPassword"`
	ApiKey   string
	BasePath string `yaml:"BasePath"`
}

type Cmdb struct {
	Config Config
	Cache  *cache.Cache
	Client *client.Client
	Logger *log.Logger
	Error  error
}

type ErrorReturn struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
}

type CiRelationDirection string

const (
	CI_RELATION_DIRECTION_ALL             CiRelationDirection = "all"
	CI_RELATION_DIRECTION_DIRECTED_FROM                       = "directed_from"
	CI_RELATION_DIRECTION_DIRECTED_TO                         = "directed_to"
	CI_RELATION_DIRECTION_BIDIRECTIONAL                       = "bidirectional"
	CI_RELATION_DIRECTION_OMNIDIRECTIONAL                     = "omnidirectional"
)

type AttributeValueType string

const (
	ATTRIBUTE_VALUE_TYPE_TEXT    AttributeValueType = "value_text"
	ATTRIBUTE_VALUE_TYPE_DATE                       = "value_date"
	ATTRIBUTE_VALUE_TYPE_DEFAULT                    = "value_default"
	ATTRIBUTE_VALUE_TYPE_CI                         = "value_ci"
)

type UpdateMode string

const (
	UPDATE_MODE_INSERT UpdateMode = "insert"
	UPDATE_MODE_UPDATE            = "update"
	UPDATE_MODE_DELETE            = "delete"
	UPDATE_MODE_SET               = "set"
)

func init() {
	log.SetLevel(log.InfoLevel)
	if os.Getenv("WORKFLOW_DEBUGGING") == "true" {
		log.SetLevel(log.DebugLevel)
	}
}

func (i *Cmdb) LoadConfigFile(configFile string) (err error) {
	_, err = os.Stat(configFile)
	if os.IsNotExist(err) {
		WorkflowConfigPath := os.Getenv("WORKFLOW_CONFIG_PATH")
		log.Debugf("WORKFLOW_CONFIG_PATH: %s", WorkflowConfigPath)
		configFile = filepath.Join(WorkflowConfigPath, configFile)
	} else if err != nil {
		return
	}

	log.Debugf("ConfigFile: %s", configFile)

	_, err = os.Stat(configFile)
	if err != nil {
		return
	}

	yamlFile, err := ioutil.ReadFile(configFile)
	if err != nil {
		return
	}

	err = i.LoadConfig(yamlFile)
	return
}

func (i *Cmdb) LoadConfig(config []byte) (err error) {
	if err = yaml.Unmarshal(config, &i.Config); err != nil {
		return
	}

	err = i.applyUrlFromRedirect()
	if err != nil {
		return
	}

	i.Client = client.New(i.Config.Url)
	return
}

func (i *Cmdb) applyUrlFromRedirect() (err error) {
	c := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	resp, err := c.Get(i.Config.Url)
	if err != nil {
		return
	}

	for name, header := range resp.Header {
		if name == "Location" && len(header) > 0 {
			baseUrl, _ := url.Parse(header[0])
			i.Config.Url = baseUrl.Scheme + "://" + baseUrl.Host + "/"
		}
	}

	return
}

// New returns a new Cmdb Client to access the V2 Api
func New() (i *Cmdb) {
	return &Cmdb{
		Cache:  cache.New(5*time.Minute, 10*time.Minute),
		Client: &client.Client{},
		Logger: log.New(),
	}
}
