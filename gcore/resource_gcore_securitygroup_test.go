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

	ipTemplate1 := fmt.Sprintf(`
			resource "gcore_securitygroup" "acctest" {
			  %s
              %s
			  name = "test"
			  metadata_map = {
				key1 = "val1"
				key2 = "val2"
			  }
			  security_group_rules {
			  	direction = "egress"
			    ethertype = "IPv4"
				protocol = "vrrp"
			  }
			}
		`, projectInfo(), regionInfo())

	ipTemplate2 := fmt.Sprintf(`
			resource "gcore_securitygroup" "acctest" {
			  %s
              %s
			  name = "test"
			  metadata_map = {
				key3 = "val3"
			  }
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
				Config: ipTemplate1,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(fullName),
					resource.TestCheckResourceAttr(fullName, "metadata_map.key1", "val1"),
					resource.TestCheckResourceAttr(fullName, "metadata_map.key2", "val2"),
					testAccCheckMetadata(fullName, true, map[string]interface{}{
						"key1": "val1",
						"key2": "val2",
					}),
				),
			},
			{
				Config: ipTemplate2,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(fullName),
					resource.TestCheckResourceAttr(fullName, "metadata_map.key3", "val3"),
					testAccCheckMetadata(fullName, true, map[string]interface{}{
						"key3": "val3",
					}),
					testAccCheckMetadata(fullName, false, map[string]interface{}{
						"key1": "val1",
					}),
					testAccCheckMetadata(fullName, false, map[string]interface{}{
						"key2": "val2",
					}),
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
