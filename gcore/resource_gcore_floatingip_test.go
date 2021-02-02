package gcore

import (
	"fmt"
	"testing"

	"github.com/G-Core/gcorelabscloud-go/gcore/floatingip/v1/floatingips"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccFloatingIP(t *testing.T) {
	fullName := "gcore_floatingip.acctest"

	ipTemplate := fmt.Sprintf(`
			resource "gcore_floatingip" "acctest" {
			  %s
              %s
			}
		`, projectInfo(), regionInfo())

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccFloatingIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: ipTemplate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(fullName),
					resource.TestCheckResourceAttr(fullName, "fixed_ip_address", ""),
					resource.TestCheckResourceAttr(fullName, "port_id", ""),
				),
			},
		},
	})
}

func testAccFloatingIPDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	client, err := CreateTestClient(config.Provider, floatingIPsPoint)
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gcore_floatingip" {
			continue
		}

		_, err := floatingips.Get(client, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("ReservedFixedIP still exists")
		}
	}

	return nil
}
