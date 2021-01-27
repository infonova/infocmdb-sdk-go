package infocmdb

import (
	"errors"
	"github.com/infonova/infocmdb-sdk-go/infocmdb/config"
	"time"

	"github.com/patrickmn/go-cache"
	log "github.com/sirupsen/logrus"
)

var (
	ErrFailedToCreateInfoCMDB  = errors.New("failed to create infocmdb object")
	ErrNoCredentials           = errors.New("must provide credentials")
	ErrNotImplemented          = errors.New("not implemented")
	ErrNoResult                = errors.New("query returned no result")
	ErrTooManyResults          = errors.New("query returned to many results, expected one")
	ErrWebserviceResponseNotOk = errors.New("webservice response was not ok")
)

type Config struct {
	ApiUrl       string `yaml:"apiUrl"`
	ApiUser      string `yaml:"apiUser"`
	ApiPassword  string `yaml:"apiPassword"`
	ApiKey       string
	CmdbBasePath string `yaml:"CmdbBasePath"`
}

type Cmdb struct {
	Config Config
	Cache  *cache.Cache
}

type CiRelationDirection string

const (
	CI_RELATION_DIRECTION_ALL             CiRelationDirection = "all"
	CI_RELATION_DIRECTION_DIRECTED_FROM   CiRelationDirection = "directed_from"
	CI_RELATION_DIRECTION_DIRECTED_TO     CiRelationDirection = "directed_to"
	CI_RELATION_DIRECTION_BIDIRECTIONAL   CiRelationDirection = "bidirectional"
	CI_RELATION_DIRECTION_OMNIDIRECTIONAL CiRelationDirection = "omnidirectional"
)

type AttributeValueType string

const (
	ATTRIBUTE_VALUE_TYPE_TEXT    AttributeValueType = "value_text"
	ATTRIBUTE_VALUE_TYPE_DATE    AttributeValueType = "value_date"
	ATTRIBUTE_VALUE_TYPE_DEFAULT AttributeValueType = "value_default"
	ATTRIBUTE_VALUE_TYPE_CI      AttributeValueType = "value_ci"
)

// New returns a new Cmdb Client to access the V1 Api
func New() (i *Cmdb) {
	return &Cmdb{
		Cache: cache.New(5*time.Minute, 10*time.Minute),
	}
}

func (i *Cmdb) LoadConfig(config Config) {
	i.Config = config
}

func (i *Cmdb) LoadConfigFile(path string) (err error) {
	err = config.LoadYamlConfig(path, &i.Config)
	if err != nil {
		return err
	}

	return
}

func (i *Cmdb) Login() error {
	if i.Config.ApiKey != "" {
		log.Trace("already logged in")
		return nil
	}

	if i.Config.ApiUser == "" {
		return ErrNoCredentials
	}
	return i.LoginWithUserPass(i.Config.ApiUrl, i.Config.ApiUser, i.Config.ApiPassword)
}
