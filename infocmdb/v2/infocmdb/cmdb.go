package infocmdb

import (
	"errors"
	"github.com/infonova/infocmdb-sdk-go/infocmdb/config"
	"net/http"
	"net/url"
	"time"

	"github.com/infonova/infocmdb-sdk-go/infocmdb/v2/infocmdb/client"
	"github.com/patrickmn/go-cache"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	Url      string `yaml:"apiUrl"`
	Username string `yaml:"apiUser"`
	Password string `yaml:"apiPassword"`
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

var (
	ErrFailedToCreateInfoCMDB  = errors.New("failed to create infocmdb object")
	ErrNoCredentials           = errors.New("must provide credentials")
	ErrNotImplemented          = errors.New("not implemented")
	ErrNoResult                = errors.New("query returned no result")
	ErrTooManyResults          = errors.New("query returned to many results, expected one")
	ErrWebserviceResponseNotOk = errors.New("webservice response was not ok")
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

func (cmdb *Cmdb) LoadConfig(config Config) {
	cmdb.Config = config
	cmdb.Client = client.New(config.Url)
	return
}

func (cmdb *Cmdb) LoadConfigFile(path string) (err error) {
	err = config.LoadYamlConfig(path, &cmdb.Config)
	if err != nil {
		return
	}

	err = cmdb.applyUrlFromRedirect()
	if err != nil {
		return
	}

	log.Debugf("Config after applied url from redirect: %+v", cmdb.Config)
	cmdb.Client = client.New(cmdb.Config.Url)
	return
}

func (cmdb *Cmdb) applyUrlFromRedirect() (err error) {
	c := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	resp, err := c.Get(cmdb.Config.Url)
	if err != nil {
		return
	}

	for name, header := range resp.Header {
		if name == "Location" && len(header) > 0 {
			baseUrl, _ := url.Parse(header[0])
			cmdb.Config.Url = baseUrl.Scheme + "://" + baseUrl.Host + "/"
		}
	}

	return
}

// New returns a new Cmdb Client to access the V2 Api
func New() (cmdb *Cmdb) {
	return &Cmdb{
		Cache:  cache.New(5*time.Minute, 10*time.Minute),
		Client: &client.Client{},
		Logger: log.New(),
	}
}
