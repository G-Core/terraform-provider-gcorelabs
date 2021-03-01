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
		Name string
	}

	create := Params{"test"}

	update := Params{"test1"}

	fullName := "gcore_loadbalancer.acctest"

	ripTemplate := func(params *Params) string {
		return fmt.Sprintf(`
			resource "gcore_loadbalancer" "acctest" {
			  %s
              %s
			  name = "%s"
			  flavor = "lb1-1-2"
              listener {
                name = "test"
                protocol = "HTTP"
                protocol_port = 80
              }
			}
		`, projectInfo(), regionInfo(), params.Name)
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
				),
			},
			{
				Config: ripTemplate(&update),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(fullName),
					resource.TestCheckResourceAttr(fullName, "name", update.Name),
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
