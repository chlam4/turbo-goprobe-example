package probe

import  (
	"fmt"

	// Turbo sdk imports
	"github.com/turbonomic/turbo-go-sdk/pkg/proto"
	"github.com/turbonomic/turbo-go-sdk/pkg/builder"
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
	return "targetIdentifier"
}

func (myProbe *ExampleRegistrationClient) GetAccountDefinition() []*proto.AccountDefEntry {
	var acctDefProps []*proto.AccountDefEntry

	// target id
	targetIDAcctDefEntry := builder.NewAccountDefEntryBuilder("targetIdentifier", "Address",
		"IP address of the probe", ".*",
		true, false).Create()

	acctDefProps = append(acctDefProps, targetIDAcctDefEntry)

	// username
	usernameAcctDefEntry := builder.NewAccountDefEntryBuilder("username", "Username",
		"Username of the probe", ".*",
		true, false).Create()
	acctDefProps = append(acctDefProps, usernameAcctDefEntry)

	// password
	passwdAcctDefEntry := builder.NewAccountDefEntryBuilder("password", "Password",
		"Password of the probe", ".*",
		true, true).Create()
	acctDefProps = append(acctDefProps, passwdAcctDefEntry)

	return acctDefProps
}