package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAzureRMSubnet_basic(t *testing.T) {

	ri := acctest.RandInt()
	config := fmt.Sprintf(testAccAzureRMSubnet_basic, ri, ri, ri, ri, ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMSubnetDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSubnetExists("azurerm_subnet.test"),
				),
			},
		},
	})
}

func TestAccAzureRMSubnet_disappears(t *testing.T) {

	ri := acctest.RandInt()
	config := fmt.Sprintf(testAccAzureRMSubnet_basic, ri, ri, ri, ri, ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMSubnetDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSubnetExists("azurerm_subnet.test"),
					testCheckAzureRMSubnetDisappears("azurerm_subnet.test"),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testCheckAzureRMSubnetExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		name := rs.Primary.Attributes["name"]
		vnetName := rs.Primary.Attributes["virtual_network_name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for subnet: %s", name)
		}

		conn := testAccProvider.Meta().(*ArmClient).subnetClient

		resp, err := conn.Get(resourceGroup, vnetName, name, "")
		if err != nil {
			return fmt.Errorf("Bad: Get on subnetClient: %s", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Subnet %q (resource group: %q) does not exist", name, resourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMSubnetDisappears(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		name := rs.Primary.Attributes["name"]
		vnetName := rs.Primary.Attributes["virtual_network_name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for subnet: %s", name)
		}

		conn := testAccProvider.Meta().(*ArmClient).subnetClient

		_, err := conn.Delete(resourceGroup, vnetName, name, make(chan struct{}))
		if err != nil {
			return fmt.Errorf("Bad: Delete on subnetClient: %s", err)
		}

		return nil
	}
}

func testCheckAzureRMSubnetDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*ArmClient).subnetClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_subnet" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		vnetName := rs.Primary.Attributes["virtual_network_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := conn.Get(resourceGroup, vnetName, name, "")

		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Subnet still exists:\n%#v", resp.SubnetPropertiesFormat)
		}
	}

	return nil
}

var testAccAzureRMSubnet_basic = `
resource "azurerm_resource_group" "test" {
    name = "acctestRG-%d"
    location = "West US"
}

resource "azurerm_virtual_network" "test" {
    name = "acctestvirtnet%d"
    address_space = ["10.0.0.0/16"]
    location = "West US"
    resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
    name = "acctestsubnet%d"
    resource_group_name = "${azurerm_resource_group.test.name}"
    virtual_network_name = "${azurerm_virtual_network.test.name}"
    address_prefix = "10.0.2.0/24"
	route_table_id = "${azurerm_route_table.test.id}" 
}

resource "azurerm_route_table" "test" {
	name = "acctestroutetable%d"
	resource_group_name = "${azurerm_resource_group.test.name}"
	location = "West US"
}

resource "azurerm_route" "test" {
	name = "acctestroute%d"
	resource_group_name = "${azurerm_resource_group.test.name}"
	route_table_name  = "${azurerm_route_table.test.name}" 

	address_prefix = "10.100.0.0/14" 
	next_hop_type = "VirtualAppliance" 
	next_hop_in_ip_address = "10.10.1.1" 
}
`
