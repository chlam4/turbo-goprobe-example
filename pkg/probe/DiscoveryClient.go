package probe

import (
	"fmt"

	// Turbo sdk imports
	"github.com/turbonomic/turbo-go-sdk/pkg/proto"
	"github.com/turbonomic/turbo-go-sdk/pkg/probe"
	"github.com/turbonomic/turbo-go-sdk/pkg/builder"

	// Example probe
	"github.com/chlam4/turbo-goprobe-example/pkg/conf"
)


// Discovery Client for the Example Probe
// Implements the TurboDiscoveryClient interface
type ExampleDiscoveryClient struct {
	clientConf        *conf.ExampleTargetConf
	targetIdentifier  string
	username          string
	pwd               string
	topoSource 	*TopologyGenerator
}

func NewDiscoveryClient(targetIdentifier string, confFile string) (*ExampleDiscoveryClient, error) {
	// Parse conf file to create clientConf
	clientConf, _ := conf.NewExampleTargetConf(confFile)
	fmt.Printf("[ExampleDiscoveryClient] Target Conf %v\n", clientConf)
	topologyAccessor, err := NewTopologyGenerator(2, 3)
	if (err != nil) {
		fmt.Errorf("Error when instantiating a topology generator", err)
		return nil, err
	}
	client := &ExampleDiscoveryClient{
		targetIdentifier: targetIdentifier,
		clientConf: clientConf,
		topoSource: topologyAccessor,
	}

	return client, nil
}


// Get the Account Values to create VMTTarget in the turbo server corresponding to this client
func (handler *ExampleDiscoveryClient) GetAccountValues() *probe.TurboTargetInfo {
	var accountValues []*proto.AccountValue
	// Convert all parameters in clientConf to AccountValue list
	prop := "Address"
	accVal := &proto.AccountValue{
		Key: &prop,
		StringValue: &handler.clientConf.Address,
	}
	accountValues = append(accountValues, accVal)

	prop = "Username"
	accVal = &proto.AccountValue{
		Key: &prop,
		StringValue: &handler.clientConf.Username,
	}
	accountValues = append(accountValues, accVal)

	prop = "Password"
	accVal = &proto.AccountValue{
		Key: &prop,
		StringValue: &handler.clientConf.Password,
	}
	accountValues = append(accountValues, accVal)

	targetInfo := probe.NewTurboTargetInfoBuilder("example", "example", "id", accountValues).Create()
	return targetInfo
}

// Validate the Target
func (handler *ExampleDiscoveryClient) Validate(accountValues[] *proto.AccountValue) (*proto.ValidationResponse, error) {
	fmt.Printf("[ExampleDiscoveryClient] BEGIN Validation for ExampleDiscoveryClient  %s", accountValues)
	// TODO: connect to the client and get validation response
	validationResponse := &proto.ValidationResponse{}

	fmt.Printf("[ExampleDiscoveryClient] validation response %s\n", validationResponse)
	return validationResponse, nil
}

// Discover the Target Topology
func (handler *ExampleDiscoveryClient) Discover(accountValues[] *proto.AccountValue) (*proto.DiscoveryResponse, error) {
	fmt.Printf("[ExampleProbe] ========= Discovery for ExampleProbe ============= %s", accountValues)
	discoveryResults, err := handler.Discover_Old()
	// 4. Build discovery response.
	// If there is error during discovery, return an ErrorDTO.
	var discoveryResponse *proto.DiscoveryResponse
	if err != nil {
		// If there is error during discovery, return an ErrorDTO.
		serverity := proto.ErrorDTO_CRITICAL
		description := fmt.Sprintf("%v", err)
		errorDTO := &proto.ErrorDTO{
			Severity:    &serverity,
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
	fmt.Printf("[ExampleProbe] discovery response %s\n", discoveryResponse)
	return discoveryResponse, nil
}


// ======================================================================
func (this *ExampleDiscoveryClient) Discover_Old() ([]*proto.EntityDTO, error) {
	var discoveryResults []*proto.EntityDTO

	pmDTOs, err := this.discoverPMs()
	if err != nil {
		return nil, fmt.Errorf("Error found during PM discovery: %s", err)
	}
	discoveryResults = append(discoveryResults, pmDTOs...)

	vmDTOs, err := this.discoverVMs()
	if err != nil {
		return nil, fmt.Errorf("Error found during VM discovery: %s", err)
	}
	discoveryResults = append(discoveryResults, vmDTOs...)

	return discoveryResults, nil
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
	var commoditiesSold []*proto.CommodityDTO
	pmResourceStat := pm.ResourceStat

	cpuComm, _ := builder.NewCommodityDTOBuilder(proto.CommodityDTO_CPU).
		Capacity(pmResourceStat.cpuCapacity).
		Used(pmResourceStat.cpuUsed).
		Create()
	commoditiesSold = append(commoditiesSold, cpuComm)

	memComm, _ := builder.NewCommodityDTOBuilder(proto.CommodityDTO_MEM).
		Capacity(pmResourceStat.memCapacity).
		Used(pmResourceStat.memUsed).
		Create()
	commoditiesSold = append(commoditiesSold, memComm)

	return commoditiesSold
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
	var commoditiesSold []*proto.CommodityDTO
	vmResourceStat := vm.ResourceStat

	vCpuComm, _ := builder.NewCommodityDTOBuilder(proto.CommodityDTO_VCPU).
		Capacity(vmResourceStat.vCpuCapacity).
		Used(vmResourceStat.vCpuUsed).
		Create()
	commoditiesSold = append(commoditiesSold, vCpuComm)

	vMemComm, _ := builder.NewCommodityDTOBuilder(proto.CommodityDTO_VMEM).
		Capacity(vmResourceStat.vMemCapacity).
		Used(vmResourceStat.vMemUsed).
		Create()
	commoditiesSold = append(commoditiesSold, vMemComm)

	return commoditiesSold
}

func createVMCommoditiesBought(vm *VirtualMachine) []*proto.CommodityDTO {
	var commoditiesBought []*proto.CommodityDTO
	vCpuCommBought, _ := builder.NewCommodityDTOBuilder(proto.CommodityDTO_CPU).
		Used(vm.ResourceStat.vCpuUsed).
		Create()
	commoditiesBought = append(commoditiesBought, vCpuCommBought)

	vMemCommBought, _ := builder.NewCommodityDTOBuilder(proto.CommodityDTO_MEM).
		Used(vm.ResourceStat.vMemCapacity).
		Create()
	commoditiesBought = append(commoditiesBought, vMemCommBought)

	return commoditiesBought
}