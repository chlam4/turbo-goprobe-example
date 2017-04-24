package main

import (
	"flag"
	"os"
	"github.com/golang/glog"

	// Turbo sdk imports
	"github.com/turbonomic/turbo-go-sdk/pkg/probe"
	"github.com/turbonomic/turbo-go-sdk/pkg/service"

	// Example probe
	example "github.com/chlam4/turbo-goprobe-example/pkg/probe"
)

func main() {
	flag.Parse()

	targetConf := "target-conf.json"
	turboCommConf := "turbo-server-conf.json"

	communicator, err := service.ParseTurboCommunicationConfig(turboCommConf)
	if err != nil {
		glog.Infof("Error while parsing the turbo communicator config file %v: %v\n", turboCommConf, err)
		os.Exit(1)
	}

	// Example Probe Registration Client
	registrationClient := &example.ExampleRegistrationClient{}
	// Example Probe Registration Client
	discoveryClient, err := example.NewDiscoveryClient(targetConf)
	if err != nil {
		glog.Infof("Error while instantiating a discovery client at %v with config %v: %v\n", turboCommConf, targetConf, err)
		os.Exit(1)
	}

	tapService, err := service.NewTAPServiceBuilder().
		WithTurboCommunicator(communicator).
		WithTurboProbe(probe.NewProbeBuilder(discoveryClient.ClientConf.TargetType, discoveryClient.ClientConf.ProbeCategory).
			RegisteredBy(registrationClient).
			DiscoversTarget(discoveryClient.ClientConf.Address, discoveryClient)).Create()

	if err != nil {
		glog.Infof("Error while building turbo tap service on target %v: %v\n", discoveryClient.ClientConf.Address, err)
		os.Exit(1)
	}

	// Connect to the Turbo server
	tapService.ConnectToTurbo()

	select {}
}
