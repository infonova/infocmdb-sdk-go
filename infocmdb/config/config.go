package config

import (
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
)

// Loads a workflow configuration file.
//
// If the given path is a absolute path to an existing file,
// it will be read and parsed in yaml format into the given config output parameter.
// If it is a relative path instead, the config path will be resolved using the WORKFLOW_CONFIG_PATH env variable.
func LoadYamlConfig(path string, config interface{}) (err error) {
	path, err = resolveAbsoluteConfigFilePath(path)
	if err != nil {
		return
	}

	return parseYamlConfig(path, config)
}

func resolveAbsoluteConfigFilePath(path string) (string, error) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		workflowConfigPath := os.Getenv("WORKFLOW_CONFIG_PATH")
		path = filepath.Join(workflowConfigPath, path)
	} else if err != nil {
		return "", err
	}

	log.Debugf("Loading workflow config file: %s", path)

	_, err = os.Stat(path)
	return path, err
}

func parseYamlConfig(path string, config interface{}) (err error) {
	configBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}

	log.Tracef("Config file content:\n%s", configBytes)

	err = yaml.Unmarshal(configBytes, config)
	if err != nil {
		return
	}

	log.Debugf("Config: %+v", config)
	return
}
