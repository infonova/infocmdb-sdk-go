package infocmdblibrary

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
)

var (
	ErrNoCredentials = errors.New("must provide credentials")
)

type Credentials struct {
	ApiKey   string
	Username string
	Password string
}

type Config struct {
	ApiUrl       string `yaml:"apiUrl"`
	ApiUser      string `yaml:"apiUser"`
	ApiPassword  string `yaml:"apiPassword"`
	CmdbBasePath string `yaml:"CmdbBasePath"`
	URL          string
}

type InfoCMDB struct {
	WS     Webservice
	WC     CmdbWebClient
	Config Config
}

func init() {
	log.SetLevel(log.InfoLevel)
	if os.Getenv("INFOCMDB_WORKFLOW_DEBUGGING") == "true" {
		log.SetLevel(log.DebugLevel)
	}
}

func (i *InfoCMDB) LoadConfig(config string) (err error) {
	InfoCmdbBasePath := filepath.Base(os.Getenv("INFOCMDB_PATH"))
	configFile := filepath.Join(InfoCmdbBasePath, config)

	_, err = os.Stat(configFile)
	if err != nil {
		return err
	}

	yamlFile, err := ioutil.ReadFile(configFile)
	if err != nil {
		return err
	}

	c := Config{}
	err = yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		return err
	}

	return nil
}

func NewCMDB(url string, cred Credentials) (i InfoCMDB, err error) {
	i.WS.client = &i.WC
	i.WC.Url = url
	i.WC.ApiKey = cred.ApiKey

	err = i.Login(url, cred)
	return i, err
}

func (i *InfoCMDB) Login(url string, cred Credentials) error {
	if i.WC.ApiKey != "" {
		log.Debug("already logged in")
		return nil
	}
	if cred.Username == "" {
		return ErrNoCredentials
	}
	return i.WC.Login(url, cred.Username, cred.Password)
}
