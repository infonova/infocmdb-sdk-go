package cmdb

import (
	"github.com/infonova/infocmdb-lib-go/core/v2/cmdb/client"
	"github.com/patrickmn/go-cache"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

type Config struct {
	Url      string `yaml:"apiUrl"`
	Username string `yaml:"apiUser"`
	Password string `yaml:"apiPassword"`
	ApiKey   string
	BasePath string `yaml:"BasePath"`
}

type InfoCMDB struct {
	Config Config
	Cache  *cache.Cache
	Client *client.Client
	Logger *log.Logger
	Error  error
}

type UpdateMode string

const (
	UpdateModeInsert UpdateMode = "insert"
	UpdateModeUpdate UpdateMode = "update"
	UpdateModeDelete UpdateMode = "delete"
	UpdateModeSet    UpdateMode = "set"
)

func init() {
	log.SetLevel(log.InfoLevel)
	if os.Getenv("WORKFLOW_DEBUGGING") == "true" {
		log.SetLevel(log.DebugLevel)
	}
}

func (i *InfoCMDB) LoadConfigFile(configFile string) *InfoCMDB {
	_, err := os.Stat(configFile)
	if os.IsNotExist(err) {
		WorkflowConfigPath := filepath.Dir(os.Getenv("WORKFLOW_CONFIG_PATH") + "/")
		log.Debugf("WORKFLOW_CONFIG_PATH: %s", WorkflowConfigPath)
		configFile = filepath.Join(WorkflowConfigPath, configFile)
	} else if err != nil {
		i.AddError(err)
		return i
	}

	log.Debugf("ConfigFile: %s", configFile)

	_, err = os.Stat(configFile)
	if err != nil {
		i.AddError(err)
		return i
	}

	yamlFile, err := ioutil.ReadFile(configFile)
	if err != nil {
		i.AddError(err)
		return i
	}

	return i.LoadConfig(yamlFile)
}

func (i *InfoCMDB) LoadConfig(config []byte) *InfoCMDB {
	if err := yaml.Unmarshal(config, &i.Config); err != nil {
		i.AddError(err)
		return i
	}

	err := i.applyUrlFromRedirect()
	if err != nil {
		i.AddError(err)
	}

	i.Client = client.NewClient(i.Config.Url)
	return i
}

func (i *InfoCMDB) applyUrlFromRedirect() (err error) {
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

func NewCMDB() (i *InfoCMDB) {
	return &InfoCMDB{
		Cache:  cache.New(5*time.Minute, 10*time.Minute),
		Client: &client.Client{},
		Logger: log.New(),
	}
}
