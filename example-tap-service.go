package main

import (
	"fmt"
	"os"

	// Turbo sdk imports
	"github.com/turbonomic/turbo-go-sdk/pkg/probe"
	"github.com/turbonomic/turbo-go-sdk/pkg/service"

	// Example probe
	example "github.com/chlam4/turbo-goprobe-example/pkg/probe"
)

func main() {
	targetConf 	:= "target-conf.json"
	targetType 	:= "ExampleGoProbe"
	probeCategory 	:= "CloudNative"

	turboCommConf 	:= "turbo-server-conf.json"
	target1 	:= "Test1"

	communicator, err := service.ParseTurboCommunicationConfig(turboCommConf)
	if err != nil {
		fmt.Printf("Error while parsing the turbo communicator config file %v: %v\n", turboCommConf, err)
		os.Exit(1)
	}

	// Example Probe Registration Client
	registrationClient := &example.ExampleRegistrationClient{}
	// Example Probe Registration Client
	discoveryClient, err := example.NewDiscoveryClient(target1, targetConf)
	if err != nil {
		fmt.Printf("Error while instantiating a discovery client at %v with config %v: %v\n", turboCommConf, targetConf, err)
		os.Exit(1)
	}

	tapService, err := service.NewTAPServiceBuilder().
		WithTurboCommunicator(communicator).
		WithTurboProbe(probe.NewProbeBuilder(targetType, probeCategory).
		RegisteredBy(registrationClient).
		DiscoversTarget(target1, discoveryClient)).Create()

	if err != nil {
		fmt.Printf("Error while building turbo tap service on target %v: %v\n", target1, err)
		os.Exit(1)
	}

	// Connect to the Turbo server
	tapService.ConnectToTurbo()

	select {}
}
