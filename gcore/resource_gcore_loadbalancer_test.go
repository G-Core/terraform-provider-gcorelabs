//go:build cloud
// +build cloud

package gcore

import (
	"fmt"
	"testing"

	"github.com/G-Core/gcorelabscloud-go/gcore/loadbalancer/v1/loadbalancers"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccLoadBalancer(t *testing.T) {
	type Params struct {
		Name        string
		MetadataMap string
	}

	create := Params{"test", `{
					key1 = "val1"
					key2 = "val2"
				}`}

	update := Params{"test1", `{
					key3 = "val3"
				}`}

	fullName := "gcore_loadbalancerv2.acctest"

	ripTemplate := func(params *Params) string {
		return fmt.Sprintf(`
			resource "gcore_loadbalancerv2" "acctest" {
			  %s
              %s
			  name = "%s"
			  flavor = "lb1-1-2"
    		  metadata_map = %s
			}
		`, projectInfo(), regionInfo(), params.Name, params.MetadataMap)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: ripTemplate(&create),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(fullName),
					resource.TestCheckResourceAttr(fullName, "name", create.Name),
					testAccCheckMetadata(fullName, true, map[string]string{
						"key1": "val1",
						"key2": "val2",
					}),
				),
			},
			{
				Config: ripTemplate(&update),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(fullName),
					resource.TestCheckResourceAttr(fullName, "name", update.Name),
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

func testAccLoadBalancerDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	client, err := CreateTestClient(config.Provider, LoadBalancersPoint, versionPointV1)
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gcore_loadbalancer" {
			continue
		}

		_, err := loadbalancers.Get(client, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("LoadBalancer still exists")
		}
	}

	return nil
}
