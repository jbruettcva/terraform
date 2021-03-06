package azurerm

import (
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/arm/network"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAzureRMLoadBalancerNatRule_basic(t *testing.T) {
	var lb network.LoadBalancer
	ri := acctest.RandInt()
	natRuleName := fmt.Sprintf("NatRule-%d", ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLoadBalancerNatRule_basic(ri, natRuleName),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLoadBalancerExists("azurerm_lb.test", &lb),
					testCheckAzureRMLoadBalancerNatRuleExists(natRuleName, &lb),
				),
			},
		},
	})
}

func TestAccAzureRMLoadBalancerNatRule_removal(t *testing.T) {
	var lb network.LoadBalancer
	ri := acctest.RandInt()
	natRuleName := fmt.Sprintf("NatRule-%d", ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMLoadBalancerNatRule_basic(ri, natRuleName),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLoadBalancerExists("azurerm_lb.test", &lb),
					testCheckAzureRMLoadBalancerNatRuleExists(natRuleName, &lb),
				),
			},
			{
				Config: testAccAzureRMLoadBalancerNatRule_removal(ri),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMLoadBalancerExists("azurerm_lb.test", &lb),
					testCheckAzureRMLoadBalancerNatRuleNotExists(natRuleName, &lb),
				),
			},
		},
	})
}

func testCheckAzureRMLoadBalancerNatRuleExists(natRuleName string, lb *network.LoadBalancer) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, _, exists := findLoadBalancerNatRuleByName(lb, natRuleName)
		if !exists {
			return fmt.Errorf("A NAT Rule with name %q cannot be found.", natRuleName)
		}

		return nil
	}
}

func testCheckAzureRMLoadBalancerNatRuleNotExists(natRuleName string, lb *network.LoadBalancer) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, _, exists := findLoadBalancerNatRuleByName(lb, natRuleName)
		if exists {
			return fmt.Errorf("A NAT Rule with name %q has been found.", natRuleName)
		}

		return nil
	}
}

func testAccAzureRMLoadBalancerNatRule_basic(rInt int, natRuleName string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name = "acctestrg-%d"
    location = "West US"
}

resource "azurerm_public_ip" "test" {
    name = "test-ip-%d"
    location = "West US"
    resource_group_name = "${azurerm_resource_group.test.name}"
    public_ip_address_allocation = "static"
}

resource "azurerm_lb" "test" {
    name = "arm-test-loadbalancer-%d"
    location = "West US"
    resource_group_name = "${azurerm_resource_group.test.name}"

    frontend_ip_configuration {
      name = "one-%d"
      public_ip_address_id = "${azurerm_public_ip.test.id}"
    }
}

resource "azurerm_lb_nat_rule" "test" {
  location = "West US"
  resource_group_name = "${azurerm_resource_group.test.name}"
  loadbalancer_id = "${azurerm_lb.test.id}"
  name = "%s"
  protocol = "Tcp"
  frontend_port = 3389
  backend_port = 3389
  frontend_ip_configuration_name = "one-%d"
}

`, rInt, rInt, rInt, rInt, natRuleName, rInt)
}

func testAccAzureRMLoadBalancerNatRule_removal(rInt int) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
    name = "acctestrg-%d"
    location = "West US"
}

resource "azurerm_public_ip" "test" {
    name = "test-ip-%d"
    location = "West US"
    resource_group_name = "${azurerm_resource_group.test.name}"
    public_ip_address_allocation = "static"
}

resource "azurerm_lb" "test" {
    name = "arm-test-loadbalancer-%d"
    location = "West US"
    resource_group_name = "${azurerm_resource_group.test.name}"

    frontend_ip_configuration {
      name = "one-%d"
      public_ip_address_id = "${azurerm_public_ip.test.id}"
    }
}
`, rInt, rInt, rInt, rInt)
}
