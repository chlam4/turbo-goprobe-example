package main

import (
	"fmt"

	// Turbo sdk imports
	"github.com/turbonomic/turbo-go-sdk/pkg/probe"
	"github.com/turbonomic/turbo-go-sdk/pkg/service"

	// Example probe
	example "github.com/chlam4/turbo-goprobe-example/pkg/probe"
)

func main() {
	targetConf 	:= "src/github.com/turbonomic/example-turbo/cmd/tap/target-conf.json"
	targetType 	:= "ExampleGoProbe"
	probeCategory 	:= "CloudNative"

	turboCommConf 	:= "src/github.com/turbonomic/example-turbo/cmd/tap/container-conf.json"
	target1 	:= "Test1"

	communicator, err := service.ParseTurboCommunicationConfig(turboCommConf)
	if err != nil {
		fmt.Errorf("Error when trying to parse the turbo communicator config file %v: %v", turboCommConf, err)
	}

	// Example Probe Registration Client
	registrationClient := &example.ExampleRegistrationClient{}
	// Example Probe Registration Client
	discoveryClient := example.NewDiscoveryClient(target1, targetConf)

	tapService, err := service.NewTAPServiceBuilder().
		WithTurboCommunicator(communicator).
		WithTurboProbe(probe.NewProbeBuilder(targetType, probeCategory).
		RegisteredBy(registrationClient).
		DiscoversTarget(target1, discoveryClient)).Create()

	if err != nil {
		fmt.Errorf("Error when trying to build turbo tap service on target %v: %v", target1, err)
	}

	// Connect to the Turbo server
	tapService.ConnectToTurbo()

	select {}
}
