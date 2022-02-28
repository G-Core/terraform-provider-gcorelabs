//go:build cloud
// +build cloud

package gcore

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/G-Core/gcorelabscloud-go/gcore/network/v1/networks"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccNetwork(t *testing.T) {

	type Params struct {
		Name string
		Type string
		Mtu  int
	}

	create := Params{
		Name: "create_test",
		Mtu:  1450,
		Type: "vxlan",
	}

	update_name := Params{
		Name: "update_test"}

	fullName := "gcore_network.acctest"
	importStateIDPrefix := fmt.Sprintf("%s:%s:", os.Getenv("TEST_PROJECT_ID"), os.Getenv("TEST_REGION_ID"))

	NetworkTemplate := func(params *Params) string {
		template := fmt.Sprintf(`
		resource "gcore_network" "acctest" {
			name = "%s"
			%s
			%s
		`, params.Name, regionInfo(), projectInfo())

		if params.Mtu != 0 {
			template += fmt.Sprintf("mtu = %d\n", params.Mtu)
		}
		if params.Type != "" {
			template += fmt.Sprintf("type = \"%s\"\n", params.Type)
		}

		return template + "\n}"
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccNetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config: NetworkTemplate(&create),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(fullName),
					resource.TestCheckResourceAttr(fullName, "name", create.Name),
					resource.TestCheckResourceAttr(fullName, "type", create.Type),
					resource.TestCheckResourceAttr(fullName, "mtu", strconv.Itoa(create.Mtu)),
				),
			},
			{
				Config: NetworkTemplate(&update_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(fullName),
					resource.TestCheckResourceAttr(fullName, "name", update_name.Name),
					resource.TestCheckResourceAttr(fullName, "type", create.Type),
					resource.TestCheckResourceAttr(fullName, "mtu", strconv.Itoa(create.Mtu)),
				),
			},
			{
				ImportStateIdPrefix: importStateIDPrefix,
				ResourceName:        fullName,
				ImportState:         true,
			},
		},
	})
}

func testAccNetworkDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	client, err := CreateTestClient(config.Provider, networksPoint, versionPointV1)
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gcore_network" {
			continue
		}

		_, err := networks.Get(client, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Network still exists")
		}
	}

	return nil
}
