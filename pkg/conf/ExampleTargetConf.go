package conf

import (
	"io/ioutil"
	"os"
	"encoding/json"
	"github.com/golang/glog"
)

// Configuration Parameters to connect to a Mesos Target
type ExampleTargetConf struct {
	Address      string
	Username     string
	Password     string
}

// Create a new ExamplelientConf from file. Other fields have default values and can be overrided.
func NewExampleTargetConf(targetConfigFilePath string) (*ExampleTargetConf, error) {

	glog.Infof("[ExampleTargetConf] Read configration from %s", targetConfigFilePath)
	metaConfig := readConfig(targetConfigFilePath)

	return metaConfig, nil
}

// Get the config from file.
func readConfig(path string) *ExampleTargetConf {
	file, e := ioutil.ReadFile(path)
	if e != nil {
		glog.Errorf("File error: %v\n", e)
		os.Exit(1)
	}
	glog.Info(string(file))

	var config ExampleTargetConf
	err := json.Unmarshal(file, &config)

	if err != nil {
		glog.Errorf("Unmarshall error :%v\n", err)
	}
	glog.Infof("Results: %+v\n", config)

	return &config
}
