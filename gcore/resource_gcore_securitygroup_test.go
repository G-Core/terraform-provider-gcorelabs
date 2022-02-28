//go:build cloud
// +build cloud

package gcore

import (
	"fmt"
	"testing"

	"github.com/G-Core/gcorelabscloud-go/gcore/securitygroup/v1/securitygroups"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccSecurityGroup(t *testing.T) {
	fullName := "gcore_securitygroup.acctest"

	ipTemplate := fmt.Sprintf(`
			resource "gcore_securitygroup" "acctest" {
			  %s
              %s
			  name = "test"
			  security_group_rules {
			  	direction = "egress"
			    ethertype = "IPv4"
				protocol = "vrrp"
			  }
			}
		`, projectInfo(), regionInfo())

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccSecurityGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: ipTemplate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(fullName),
				),
			},
		},
	})
}

func testAccSecurityGroupDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	client, err := CreateTestClient(config.Provider, securityGroupPoint, versionPointV1)
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gcore_securitygroup" {
			continue
		}

		_, err := securitygroups.Get(client, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("SecurityGroup still exists")
		}
	}

	return nil
}
