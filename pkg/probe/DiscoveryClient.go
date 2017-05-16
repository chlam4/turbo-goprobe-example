package probe

import (
	"fmt"
	"github.com/golang/glog"

	// Turbo sdk imports
	"github.com/turbonomic/turbo-go-sdk/pkg/builder"
	"github.com/turbonomic/turbo-go-sdk/pkg/probe"
	"github.com/turbonomic/turbo-go-sdk/pkg/proto"

	// Example probe
	"github.com/turbonomic/turbo-goprobe-example/pkg/conf"
)

// Discovery Client for the Example Probe
// Implements the TurboDiscoveryClient interface
type ExampleDiscoveryClient struct {
	ClientConf *conf.ExampleTargetConf
	topoSource *TopologyGenerator
}

func NewDiscoveryClient(confFile string) (*ExampleDiscoveryClient, error) {
	// Parse conf file to create clientConf
	clientConf, _ := conf.NewExampleTargetConf(confFile)
	glog.Infof("[ExampleDiscoveryClient] Target Conf %v\n", clientConf)
	topologyAccessor, err := NewTopologyGenerator(2, 3)
	if err != nil {
		glog.Errorf("Error when instantiating a topology generator", err)
		return nil, err
	}
	client := &ExampleDiscoveryClient{
		ClientConf: clientConf,
		topoSource: topologyAccessor,
	}

	return client, nil
}

// Get the Account Values to create VMTTarget in the turbo server corresponding to this client
func (handler *ExampleDiscoveryClient) GetAccountValues() *probe.TurboTargetInfo {
	// Convert all parameters in clientConf to AccountValue list
	clientConf := handler.ClientConf

	targetId := TargetIdField
	targetIdVal := &proto.AccountValue{
		Key:         &targetId,
		StringValue: &clientConf.Address,
	}

	username := Username
	usernameVal := &proto.AccountValue{
		Key:         &username,
		StringValue: &clientConf.Username,
	}

	password := Password
	passwordVal := &proto.AccountValue{
		Key:         &password,
		StringValue: &clientConf.Password,
	}

	accountValues := []*proto.AccountValue{
		targetIdVal,
		usernameVal,
		passwordVal,
	}

	targetInfo := probe.NewTurboTargetInfoBuilder(clientConf.ProbeCategory, clientConf.TargetType, TargetIdField, accountValues).Create()
	return targetInfo
}

// Validate the Target
func (handler *ExampleDiscoveryClient) Validate(accountValues []*proto.AccountValue) (*proto.ValidationResponse, error) {
	glog.Infof("[ExampleDiscoveryClient] BEGIN Validation for ExampleDiscoveryClient %s\n", accountValues)
	// TODO: connect to the client and get validation response
	validationResponse := &proto.ValidationResponse{}

	glog.Infof("[ExampleDiscoveryClient] validation response %s\n", validationResponse)
	return validationResponse, nil
}

// Discover the Target Topology
func (handler *ExampleDiscoveryClient) Discover(accountValues []*proto.AccountValue) (*proto.DiscoveryResponse, error) {
	glog.Infof("[ExampleProbe] ========= Discovery for ExampleProbe ============= %s\n", accountValues)
	discoveryResults, err := handler.Discover_Old()
	// 4. Build discovery response.
	// If there is error during discovery, return an ErrorDTO.
	var discoveryResponse *proto.DiscoveryResponse
	if err != nil {
		// If there is error during discovery, return an ErrorDTO.
		severity := proto.ErrorDTO_CRITICAL
		description := fmt.Sprintf("%v", err)
		errorDTO := &proto.ErrorDTO{
			Severity:    &severity,
			Description: &description,
		}
		discoveryResponse = &proto.DiscoveryResponse{
			ErrorDTO: []*proto.ErrorDTO{errorDTO},
		}
	} else {
		// No error. Return the result entityDTOs.
		discoveryResponse = &proto.DiscoveryResponse{
			EntityDTO: discoveryResults,
		}
	}
	glog.Infof("[ExampleProbe] discovery response %s\n", discoveryResponse)
	return discoveryResponse, nil
}

// ======================================================================
func (this *ExampleDiscoveryClient) Discover_Old() ([]*proto.EntityDTO, error) {

	pmDTOs, err := this.discoverPMs()
	if err != nil {
		return nil, fmt.Errorf("Error found during PM discovery: %s", err)
	}

	vmDTOs, err := this.discoverVMs()
	if err != nil {
		return nil, fmt.Errorf("Error found during VM discovery: %s", err)
	}

	return append(pmDTOs, vmDTOs...), nil
}

func (this *ExampleDiscoveryClient) discoverPMs() ([]*proto.EntityDTO, error) {
	var result []*proto.EntityDTO

	pms := this.topoSource.GetPMs()
	for _, pm := range pms {
		commoditiesSold := createPMCommoditiesSold(pm)

		entityDTO, err := builder.NewEntityDTOBuilder(proto.EntityDTO_PHYSICAL_MACHINE, pm.UUID).
			DisplayName(pm.Name).
			SellsCommodities(commoditiesSold).
			Create()
		if err != nil {
			return nil, fmt.Errorf("Error creating entityDTO for PM %s: %v", pm.Name, err)
		}
		result = append(result, entityDTO)
	}

	return result, nil
}

func createPMCommoditiesSold(pm *PhysicalMachine) []*proto.CommodityDTO {

	cpuComm, _ := builder.NewCommodityDTOBuilder(proto.CommodityDTO_CPU).
		Capacity(pm.ResourceStat.cpuCapacity).
		Used(pm.ResourceStat.cpuUsed).
		Create()

	memComm, _ := builder.NewCommodityDTOBuilder(proto.CommodityDTO_MEM).
		Capacity(pm.ResourceStat.memCapacity).
		Used(pm.ResourceStat.memUsed).
		Create()

	return []*proto.CommodityDTO{cpuComm, memComm}
}

func (this *ExampleDiscoveryClient) discoverVMs() ([]*proto.EntityDTO, error) {
	var result []*proto.EntityDTO

	vms := this.topoSource.GetVMs()
	for _, vm := range vms {
		commoditiesSold := createVMCommoditiesSold(vm)
		commoditiesBought := createVMCommoditiesBought(vm)

		entityDTO, err := builder.NewEntityDTOBuilder(proto.EntityDTO_VIRTUAL_MACHINE, vm.UUID).
			DisplayName(vm.Name).
			SellsCommodities(commoditiesSold).
			Provider(builder.CreateProvider(proto.EntityDTO_PHYSICAL_MACHINE, vm.providerID)).
			BuysCommodities(commoditiesBought).
			Create()
		if err != nil {
			return nil, fmt.Errorf("Error creating entityDTO for VM %s: %v", vm.Name, err)
		}
		result = append(result, entityDTO)
	}

	return result, nil
}

func createVMCommoditiesSold(vm *VirtualMachine) []*proto.CommodityDTO {

	vCpuComm, _ := builder.NewCommodityDTOBuilder(proto.CommodityDTO_VCPU).
		Capacity(vm.ResourceStat.vCpuCapacity).
		Used(vm.ResourceStat.vCpuUsed).
		Create()

	vMemComm, _ := builder.NewCommodityDTOBuilder(proto.CommodityDTO_VMEM).
		Capacity(vm.ResourceStat.vMemCapacity).
		Used(vm.ResourceStat.vMemUsed).
		Create()

	return []*proto.CommodityDTO{vCpuComm, vMemComm}
}

func createVMCommoditiesBought(vm *VirtualMachine) []*proto.CommodityDTO {
	vCpuCommBought, _ := builder.NewCommodityDTOBuilder(proto.CommodityDTO_CPU).
		Used(vm.ResourceStat.vCpuUsed).
		Create()

	vMemCommBought, _ := builder.NewCommodityDTOBuilder(proto.CommodityDTO_MEM).
		Used(vm.ResourceStat.vMemCapacity).
		Create()

	return []*proto.CommodityDTO{vCpuCommBought, vMemCommBought}
}
