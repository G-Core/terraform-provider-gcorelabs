package gcore

import (
	"fmt"
	"testing"

	"github.com/G-Core/gcorelabscloud-go/gcore/reservedfixedip/v1/reservedfixedips"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccReservedFixedIP(t *testing.T) {
	type Params struct {
		Type  string
		IsVip bool
	}

	createExternal := Params{
		Type:  "external",
		IsVip: true,
	}

	updateExternal := Params{
		Type:  "external",
		IsVip: false,
	}

	fullName := "gcore_reservedfixedip.acctest"

	ripTemplateExternal := func(params *Params) string {
		return fmt.Sprintf(`
			resource "gcore_reservedfixedip" "acctest" {
			  %s
              %s
			  is_vip = %t
			  type = "%s"
			}
		`, projectInfo(), regionInfo(), params.IsVip, params.Type)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccReservedFixedIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: ripTemplateExternal(&createExternal),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(fullName),
					resource.TestCheckResourceAttr(fullName, "type", createExternal.Type),
					resource.TestCheckResourceAttr(fullName, "is_vip", fmt.Sprintf("%t", createExternal.IsVip)),
				),
			},
			{
				Config: ripTemplateExternal(&updateExternal),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(fullName),
					resource.TestCheckResourceAttr(fullName, "type", updateExternal.Type),
					resource.TestCheckResourceAttr(fullName, "is_vip", fmt.Sprintf("%t", updateExternal.IsVip)),
				),
			},
		},
	})
}

func testAccReservedFixedIPDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	client, err := CreateTestClient(config.Provider, floatingIPsPoint, versionPointV1)
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gcore_reservedfixedip" {
			continue
		}

		_, err := reservedfixedips.Get(client, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("ReservedFixedIP still exists")
		}
	}

	return nil
}
