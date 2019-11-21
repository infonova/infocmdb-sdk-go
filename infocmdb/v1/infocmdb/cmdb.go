package infocmdb

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/patrickmn/go-cache"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
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

func (i *Cmdb) LoadConfigFile(configFile string) (err error) {
	_, err = os.Stat(configFile)
	if err == nil {
		log.Debugf("ConfigFile found with given string: %s", configFile)
	} else {
		WorkflowConfigPath := os.Getenv("WORKFLOW_CONFIG_PATH")
		log.Debugf("WORKFLOW_CONFIG_PATH: %s", WorkflowConfigPath)
		configFile = filepath.Join(WorkflowConfigPath, configFile)
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

	return i.LoadConfig(yamlFile)
}

func (i *Cmdb) LoadConfig(config []byte) (err error) {
	return yaml.Unmarshal(config, &i.Config)
}

func New(config string) (i *Cmdb, err error) {
	i = new(Cmdb)
	err = i.LoadConfigFile(config)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	i.Cache = cache.New(5*time.Minute, 10*time.Minute)

	return i, nil
}

func (i *Cmdb) Login() error {
	if i.Config.ApiKey != "" {
		log.Debug("already logged in")
		return nil
	}

	if i.Config.ApiUser == "" {
		return ErrNoCredentials
	}
	return i.LoginWithUserPass(i.Config.ApiUrl, i.Config.ApiUser, i.Config.ApiPassword)
}
