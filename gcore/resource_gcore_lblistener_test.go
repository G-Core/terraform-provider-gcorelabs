package gcore

import (
	"fmt"
	"testing"

	"github.com/G-Core/gcorelabscloud-go/gcore/loadbalancer/v1/listeners"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccLBListener(t *testing.T) {
	type Params struct {
		Name string
	}

	create := Params{"test"}

	update := Params{"test1"}

	fullName := "gcore_lblistener.acctest"

	ripTemplate := func(params *Params) string {
		return fmt.Sprintf(`
            resource "gcore_lblistener" "acctest" {
			  %s
              %s
			  name = "%s"
			  protocol = "TCP"
			  protocol_port = 36621
			  loadbalancer_id = "%s"
			}
		`, projectInfo(), regionInfo(), params.Name, GCORE_LB_ID)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckLBListener(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccLBListenerDestroy,
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

func testAccLBListenerDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	client, err := CreateTestClient(config.Provider, LBListenersPoint)
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "gcore_lblistener" {
			_, err := listeners.Get(client, rs.Primary.ID).Extract()
			if err == nil {
				return fmt.Errorf("LBListener still exists")
			}
		}
	}

	return nil
}
