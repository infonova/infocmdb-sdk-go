package cmdb

import (
	v1 "git.appteam.infonova.cloud/infocmdb/library/core/v1/cmdb"
	v2 "git.appteam.infonova.cloud/infocmdb/library/core/v2/cmdb"
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
