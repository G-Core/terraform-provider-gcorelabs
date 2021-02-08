package gcore

import (
	"fmt"
	"testing"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/loadbalancer/v1/lbpools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccLBMember(t *testing.T) {
	type Params struct {
		Address string
		Port    string
		Weight  string
	}

	create := Params{"10.10.2.15", "8080", "1"}

	update := Params{"10.10.2.16", "8081", "5"}

	fullName := "gcore_lbmember.acctest"

	tpl := func(params *Params) string {
		return fmt.Sprintf(`
            resource "gcore_lbmember" "acctest" {
			  %s
              %s
			  pool_id = "%s"
			  address = "%s"
			  protocol_port = %s
			  weight = %s
			}
		`, projectInfo(), regionInfo(), GCORE_LBPOOL_ID, params.Address, params.Port, params.Weight)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckLBMember(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccLBMemberDestroy,
		Steps: []resource.TestStep{
			{
				Config: tpl(&create),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(fullName),
					resource.TestCheckResourceAttr(fullName, "address", create.Address),
					resource.TestCheckResourceAttr(fullName, "protocol_port", create.Port),
					resource.TestCheckResourceAttr(fullName, "weight", create.Weight),
				),
			},
			{
				Config: tpl(&update),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(fullName),
					resource.TestCheckResourceAttr(fullName, "address", update.Address),
					resource.TestCheckResourceAttr(fullName, "protocol_port", update.Port),
					resource.TestCheckResourceAttr(fullName, "weight", update.Weight),
				),
			},
		},
	})
}

func testAccLBMemberDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	client, err := CreateTestClient(config.Provider, LBPoolsPoint)
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gcore_lbmember" {
			continue
		}

		pl, err := lbpools.Get(client, GCORE_LBPOOL_ID).Extract()
		if err != nil {
			switch err.(type) {
			case gcorecloud.ErrDefault404:
				return nil
			default:
				return err
			}
		}

		for _, m := range pl.Members {
			if rs.Primary.ID == m.ID {
				return fmt.Errorf("LBMember still exists")
			}
		}
	}

	return nil
}
