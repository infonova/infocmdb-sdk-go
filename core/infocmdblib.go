package core

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
)

var (
	ErrArgumentsMissing       = errors.New("arguments missing")
	ErrFailedToCreateInfoCMDB = errors.New("failed to create infocmdb object")
	ErrNoCredentials          = errors.New("must provide credentials")
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
	if os.Getenv("INFOCMDB_WORKFLOW_DEBUGGING") == "true" {
		log.SetLevel(log.DebugLevel)
	}
}

func (i *InfoCMDB) LoadConfig(configFile string) (err error) {
	_, err = os.Stat(configFile)
	if err == nil {
		log.Debugf("Configfile found with given string: %s", configFile)
	} else {
		InfoCmdbBasePath := filepath.Dir(os.Getenv("INFOCMDB_WORKFLOW_CONFIG_PATH") + "/")
		log.Debugf("INFOCMDB_WORKFLOW_CONFIG_PATH: %s, InfoCmdbBasePath: %s", os.Getenv("INFOCMDB_WORKFLOW_CONFIG_PATH"), InfoCmdbBasePath)
		configFile = filepath.Join(InfoCmdbBasePath, configFile)
	}

	log.Debugf("Configfile: %s", configFile)

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
