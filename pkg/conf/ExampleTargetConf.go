package conf

import (
	"fmt"
	"io/ioutil"
	"os"
	"encoding/json"
)

// Configuration Parameters to connect to an example Target
type ExampleTargetConf struct {
	Address      string
	Username     string
	Password     string
}

// Create a new ExamplelientConf from file. Other fields have default values and can be overrided.
func NewExampleTargetConf(targetConfigFilePath string) (*ExampleTargetConf, error) {

	fmt.Printf("[ExampleTargetConf] Read configration from %s", targetConfigFilePath)
	metaConfig := readConfig(targetConfigFilePath)

	return metaConfig, nil
}

// Get the config from file.
func readConfig(path string) *ExampleTargetConf {
	file, e := ioutil.ReadFile(path)
	if e != nil {
		fmt.Errorf("File error: %v\n", e)
		os.Exit(1)
	}
	fmt.Print(string(file))

	var config ExampleTargetConf
	err := json.Unmarshal(file, &config)

	if err != nil {
		fmt.Errorf("Unmarshall error :%v\n", err)
	}
	fmt.Printf("Results: %+v\n", config)

	return &config
}
