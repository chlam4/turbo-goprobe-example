package conf

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"github.com/golang/glog"
)

// Configuration Parameters to connect to an example Target
type ExampleTargetConf struct {
	Address       string
	Username      string
	Password      string
	ProbeCategory string
	TargetType    string
}

// Create a new ExampleClientConf from file. Other fields have default values and can be overridden.
func NewExampleTargetConf(targetConfigFilePath string) (*ExampleTargetConf, error) {

	glog.Infof("[ExampleTargetConf] Read configuration from %s\n", targetConfigFilePath)
	metaConfig := readConfig(targetConfigFilePath)

	return metaConfig, nil
}

// Get the config from file.
func readConfig(path string) *ExampleTargetConf {
	file, e := ioutil.ReadFile(path)
	if e != nil {
		glog.Infof("File error: %v\n", e)
		os.Exit(1)
	}
	glog.Infoln(string(file))

	var config ExampleTargetConf
	err := json.Unmarshal(file, &config)

	if err != nil {
		glog.Errorf("Unmarshall error :%v\n", err)
	}
	glog.Infof("Results: %+v\n", config)

	return &config
}
