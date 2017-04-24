package probe

import (
	"github.com/turbonomic/turbo-go-sdk/pkg/proto"
	builder "github.com/turbonomic/turbo-go-sdk/pkg/supplychain"
)

var (
	cpuType  = proto.CommodityDTO_CPU
	memType  = proto.CommodityDTO_MEM
	vCpuType = proto.CommodityDTO_VCPU
	vMemType = proto.CommodityDTO_VMEM

	//Commodity key is optional, when key is set, it serves as a constraint between seller and buyer
	//for example, the buyer can only go to a seller that sells the commodity with the required key
	cpuCommKey string = "cpu_comm_key"
	memCommKey string = "mem_comm_key"

	cpuTemplateComm  *proto.TemplateCommodity = &proto.TemplateCommodity{CommodityType: &cpuType}
	memTemplateComm  *proto.TemplateCommodity = &proto.TemplateCommodity{CommodityType: &memType}
	vCpuTemplateComm *proto.TemplateCommodity = &proto.TemplateCommodity{CommodityType: &vCpuType}
	vMemTemplateComm *proto.TemplateCommodity = &proto.TemplateCommodity{CommodityType: &vMemType}
)

type SupplyChainFactory struct{}

// SupplyChain definition: this function defines the buyer/seller relationships between each of
// the entity types in * the Target, the default Supply Chain definition in this function is:
// a Virtual Machine buyer, a Physical Machine seller and the commodities are CPU and Memory.
// Each entity type and the relationships are defined by a single TemplateDTO struct
// The function returns an array of TemplateDTO pointers
// TO MODIFY:
// For each entity: Create a supply chain builder object with builder.NewSupplyChainNodeBuilder()
//		    Set a provider type if the new entity is a buyer , create commodity objects
//		    and add them to the entity's supply chain builder object
//                  Add commodity objects with the selling function to the entity you create if
//		    it is a seller.
//		    Add the new entity to the supplyChainBuilder instance with either the Top()
//		    or  Entity() methods
// The SupplyChainBuilder() function is only called once, in this function.
func (this *SupplyChainFactory) CreateSupplyChain() ([]*proto.TemplateDTO, error) {
	vmSupplyChainNodeBuilder := this.virtualMachineSupplyChainNodeBuilder()
	pmSupplyChainNodeBuilder := this.physicalMachineSupplyChainNodeBuilder()

	vmNode, _ := vmSupplyChainNodeBuilder.Create()
	pmNode, _ := pmSupplyChainNodeBuilder.Create()

	// SupplyChain building
	// The last buyer in the supply chain is set as the top entity with the Top() method
	// All other entities are added to the SupplyChainBuilder with the Entity() method
	return builder.NewSupplyChainBuilder().
		Top(vmNode).
		Entity(pmNode).
		Create()
}

// Create supply chain definition for Physical Machine.
func (this *SupplyChainFactory) physicalMachineSupplyChainNodeBuilder() *builder.SupplyChainNodeBuilder {
	// PM Creation Process
	pmSupplyChainNodeBuilder := builder.NewSupplyChainNodeBuilder(proto.EntityDTO_PHYSICAL_MACHINE)
	// Creates a Physical Machine entity and sets the type of commodity it sells to CPU
	pmSupplyChainNodeBuilder = pmSupplyChainNodeBuilder.
		Sells(cpuTemplateComm).
		Sells(memTemplateComm)

	return pmSupplyChainNodeBuilder
}

// Create supply chain definition for Vitual Machine
func (this *SupplyChainFactory) virtualMachineSupplyChainNodeBuilder() *builder.SupplyChainNodeBuilder {
	// VM Creation Process
	vmSupplyChainNodeBuilder := builder.NewSupplyChainNodeBuilder(proto.EntityDTO_VIRTUAL_MACHINE).
		Sells(vCpuTemplateComm).
		Sells(vMemTemplateComm)

	// The Entity type for the Virtual Machine's commodity provider is defined by the Provider() method.
	// The Commodity type for Virtual Machine's buying relationship is defined by the Buys() method
	vmSupplyChainNodeBuilder = vmSupplyChainNodeBuilder.
		Provider(proto.EntityDTO_PHYSICAL_MACHINE, proto.Provider_HOSTING).
		Buys(cpuTemplateComm).
		Buys(memTemplateComm)

	return vmSupplyChainNodeBuilder
}
