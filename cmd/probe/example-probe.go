package main

import (
	"flag"

	// Turbo sdk imports
	"github.com/turbonomic/turbo-go-sdk/pkg/probe"
	mediation "github.com/turbonomic/turbo-go-sdk/pkg/mediationcontainer"

	// Example probe
	example "github.com/chlam4/example-goprobe/pkg/probe"
	"github.com/golang/glog"
)

func init() {
}

func main() {
	flag.Parse()
	//switch {
	//case 1 : glog.V(3).Infof("Test")
	//}
	targetConf 	:= "src/github.com/turbonomic/example-turbo/cmd/target-conf.json"
	targetType 	:= "ExampleGoProbe"
	probeCategory 	:= "CloudNative"

	target1 	:= "Test1"

	// The WebSocket Container
	containerConfig := mediation.MediationContainerConfig{}
	mediation.CreateMediationContainer(&containerConfig)


	// Example Probe Registration Client
	registrationClient := &example.ExampleRegistrationClient{}
	// Example Probe Registration Client
	discoveryClient := example.NewDiscoveryClient(target1, targetConf)

	turboProbe, err := probe.NewProbeBuilder(targetType, probeCategory).
					RegisteredBy(registrationClient).
					DiscoversTarget(target1, discoveryClient).Create()

	if err != nil {
		glog.Errorf("Error encountered when trying to build turbo goprobe on target %v: %v", target1, err)
	}

	// Load the probe in the container
	mediation.LoadProbe(turboProbe)
	mediation.GetProbe(turboProbe.ProbeType)

	IsRegistered := make(chan bool, 1)

	// Connect to the Turbo server
	mediation.InitMediationContainer(IsRegistered)

	// Block till a message arrives on the channel
	status := <- IsRegistered
	if !status {
		glog.Infof("Probe " + turboProbe.ProbeCategory + "::" + turboProbe.ProbeType + " should be registered before adding Targets")
		return
	}
	glog.Infof("Probe " + turboProbe.ProbeCategory + "::" + turboProbe.ProbeType +" Registered : ============ Add Targets ========")

	// TODO: move inside to example probe
	//topologyAccessor, err := probe.NewTopologyGenerator(2, 3)
	//if err != nil {
	//	glog.Fatal("Error getting topology accessor: %v", err)
	//}
	//
	//stopCh := make(chan struct{})
	//defer close(stopCh)
	//go func() {
	//	for {
	//		select {
	//		case <-stopCh:
	//			return
	//		default:
	//		}
	//
	//		topologyAccessor.UpdateResource()
	//
	//		t := time.NewTimer(time.Minute * 1)
	//		select {
	//		case <-stopCh:
	//			return
	//		case <-t.C:
	//		}
	//	}
	//
	//}()

	select {}
}
