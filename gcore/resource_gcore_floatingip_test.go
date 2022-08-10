//go:build cloud
// +build cloud

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

	tpl := func(metadataMap string) string {
		return fmt.Sprintf(`
			resource "gcore_floatingip" "acctest" {
			  %s
              %s
 	          metadata_map = %s
			}
		`, projectInfo(), regionInfo(), metadataMap)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccFloatingIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: tpl(`{
					key1 = "val1"
					key2 = "val2"
				}`),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(fullName),
					resource.TestCheckResourceAttr(fullName, "fixed_ip_address", ""),
					resource.TestCheckResourceAttr(fullName, "port_id", ""),
					resource.TestCheckResourceAttr(fullName, "metadata_map.key1", "val1"),
					resource.TestCheckResourceAttr(fullName, "metadata_map.key2", "val2"),
					testAccCheckMetadata(fullName, true, map[string]string{
						"key1": "val1",
						"key2": "val2",
					}),
				),
			},
			{
				Config: tpl(`{
					key3 = "val3"
				}`),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(fullName),
					resource.TestCheckResourceAttr(fullName, "fixed_ip_address", ""),
					resource.TestCheckResourceAttr(fullName, "port_id", ""),
					testAccCheckMetadata(fullName, true, map[string]string{
						"key3": "val3",
					}),
					testAccCheckMetadata(fullName, false, map[string]string{
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

func testAccFloatingIPDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	client, err := CreateTestClient(config.Provider, floatingIPsPoint, versionPointV1)
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
