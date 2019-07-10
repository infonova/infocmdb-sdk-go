package core

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
)

var (
	ErrArgumentsMissing        = errors.New("arguments missing")
	ErrFailedToCreateInfoCMDB  = errors.New("failed to create infocmdb object")
	ErrNoCredentials           = errors.New("must provide credentials")
	ErrNotImplemented          = errors.New("not implemented")
	ErrNoResult                = errors.New("query returned no result")
	ErrTooManyResults          = errors.New("query returned to many results, expected one")
	ErrWebserviceResponseNotOk = errors.New("webservice response was not ok")
)

const (
	ATTRIBUTE_VALUE_TYPE_TEXT    = "value_text"
	ATTRIBUTE_VALUE_TYPE_DATE    = "value_date"
	ATTRIBUTE_VALUE_TYPE_DEFAULT = "value_default"
	ATTRIBUTE_VALUE_TYPE_CI      = "value_ci"

	CI_RELATION_DIRECTION_ALL             = "all"
	CI_RELATION_DIRECTION_DIRECTED_FROM   = "directed_from"
	CI_RELATION_DIRECTION_DIRECTED_TO     = "directed_to"
	CI_RELATION_DIRECTION_BIDIRECTIONAL   = "bidirectional"
	CI_RELATION_DIRECTION_OMNIDIRECTIONAL = "omnidirectional"
)

type Config struct {
	ApiUrl       string `yaml:"apiUrl"`
	ApiUser      string `yaml:"apiUser"`
	ApiPassword  string `yaml:"apiPassword"`
	ApiKey       string
	CmdbBasePath string `yaml:"CmdbBasePath"`
}

type InfoCMDB struct {
	Config Config
}

func init() {
	log.SetLevel(log.InfoLevel)
	if os.Getenv("WORKFLOW_DEBUGGING") == "true" {
		log.SetLevel(log.DebugLevel)
	}
}

func (i *InfoCMDB) LoadConfig(configFile string) (err error) {
	_, err = os.Stat(configFile)
	if err == nil {
		log.Debugf("ConfigFile found with given string: %s", configFile)
	} else {
		WorkflowConfigPath := filepath.Dir(os.Getenv("WORKFLOW_CONFIG_PATH") + "/")
		log.Debugf("WORKFLOW_CONFIG_PATH: %s", WorkflowConfigPath)
		configFile = filepath.Join(WorkflowConfigPath, configFile)
	}

	log.Debugf("ConfigFile: %s", configFile)

	_, err = os.Stat(configFile)
	if err != nil {
		return err
	}

	yamlFile, err := ioutil.ReadFile(configFile)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(yamlFile, &i.Config)
}

func NewCMDB(config string) (i *InfoCMDB, err error) {
	i = new(InfoCMDB)
	err = i.LoadConfig(config)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return i, nil
}

func (i *InfoCMDB) Login() error {
	if i.Config.ApiKey != "" {
		log.Debug("already logged in")
		return nil
	}
	if i.Config.ApiUser == "" {
		return ErrNoCredentials
	}
	return i.LoginWithUserPass(i.Config.ApiUrl, i.Config.ApiUser, i.Config.ApiPassword)
}

func (i *InfoCMDB) FunctionError(msg string) error {
	pc := make([]uintptr, 15)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()

	fullMsg := frame.Function + ": " + msg

	return errors.New(fullMsg)
}
