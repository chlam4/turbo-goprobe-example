package probe

import (
	"fmt"

	// Turbo sdk imports
	"github.com/turbonomic/turbo-go-sdk/pkg/builder"
	"github.com/turbonomic/turbo-go-sdk/pkg/proto"
)

const (
	TargetIdField string = "targetIdentifier"
	Username      string = "username"
	Password      string = "password"
)

// Registration Client for the Example Probe
// Implements the TurboRegistrationClient interface
type ExampleRegistrationClient struct {
}

func (myProbe *ExampleRegistrationClient) GetSupplyChainDefinition() []*proto.TemplateDTO {
	fmt.Println("[ExampleRegistrationClient] .......... Now use builder to create a supply chain ..........")

	// 2. Build supply chain.
	supplyChainFactory := &SupplyChainFactory{}
	templateDtos, err := supplyChainFactory.CreateSupplyChain()
	if err != nil {
		fmt.Println("[ExampleProbe] Error creating Supply chain for the example probe")
		return nil
	}
	fmt.Println("[ExampleProbe] Supply chain for the example probe is created.")
	return templateDtos
}

func (registrationClient *ExampleRegistrationClient) GetIdentifyingFields() string {
	return TargetIdField
}

func (myProbe *ExampleRegistrationClient) GetAccountDefinition() []*proto.AccountDefEntry {
	// target id
	targetIDAcctDefEntry := builder.NewAccountDefEntryBuilder(TargetIdField, "Address",
		"IP address of the target", ".*",
		true, false).Create()

	// username
	usernameAcctDefEntry := builder.NewAccountDefEntryBuilder(Username, "Username",
		"Username for the target", ".*",
		true, false).Create()

	// password
	passwdAcctDefEntry := builder.NewAccountDefEntryBuilder(Password, "Password",
		"Password for the target", ".*",
		true, true).Create()

	return []*proto.AccountDefEntry{
		targetIDAcctDefEntry,
		usernameAcctDefEntry,
		passwdAcctDefEntry,
	}
}
