package core

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

type InfoCMDB struct {
	v1 *v1.InfoCMDB
	v2 *v2.InfoCMDB
}

func init() {
	log.SetLevel(log.InfoLevel)
	if os.Getenv("WORKFLOW_DEBUGGING") == "true" {
		log.SetLevel(log.DebugLevel)
	}
}

func NewCMDB() (i *InfoCMDB, err error) {
	return new(InfoCMDB), nil
}

func (i *InfoCMDB) LoadConfig(f string) (err error) {
	if err = i.v1.LoadConfig("/app/application/configs/workflows/infocmdb.yml"); err != nil {
		return
	}
	if err = i.v2.LoadConfigFile("/app/application/configs/workflows/infocmdb.yml").Error; err != nil {
		return
	}

	return
}
