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
		Name        string
		Type        string
		Mtu         int
		MetadataMap string
	}

	paramsCreate := Params{
		Name: "create_test",
		Mtu:  1450,
		Type: "vxlan",
		MetadataMap: `{
				key1 = "val1"
				key2 = "val2"
			}`,
	}

	paramsUpdate := Params{
		Name: "update_test",
		MetadataMap: `{
				key3 = "val3"
			  }`,
	}

	fullName := "gcore_network.acctest"
	importStateIDPrefix := fmt.Sprintf("%s:%s:", os.Getenv("TEST_PROJECT_ID"), os.Getenv("TEST_REGION_ID"))

	NetworkTemplate := func(params *Params) string {
		template := fmt.Sprintf(`
		resource "gcore_network" "acctest" {
			name = "%s"
	  		metadata_map = %s
			%s
			%s
		`, params.Name, params.MetadataMap, regionInfo(), projectInfo())

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
				Config: NetworkTemplate(&paramsCreate),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(fullName),
					resource.TestCheckResourceAttr(fullName, "name", paramsCreate.Name),
					resource.TestCheckResourceAttr(fullName, "type", paramsCreate.Type),
					resource.TestCheckResourceAttr(fullName, "mtu", strconv.Itoa(paramsCreate.Mtu)),
					resource.TestCheckResourceAttr(fullName, "metadata_map.key1", "val1"),
					resource.TestCheckResourceAttr(fullName, "metadata_map.key2", "val2"),
					testAccCheckMetadata(fullName, true, map[string]string{
						"key1": "val1",
						"key2": "val2",
					}),
				),
			},
			{
				Config: NetworkTemplate(&paramsUpdate),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(fullName),
					resource.TestCheckResourceAttr(fullName, "name", paramsUpdate.Name),
					resource.TestCheckResourceAttr(fullName, "type", paramsCreate.Type),
					resource.TestCheckResourceAttr(fullName, "mtu", strconv.Itoa(paramsCreate.Mtu)),
					testAccCheckMetadata(fullName, true, map[string]string{
						"key3": "val3",
					}),
					testAccCheckMetadata(fullName, false, map[string]string{
						"key1": "val1",
					}),
					testAccCheckMetadata(fullName, false, map[string]string{
						"key2": "val2",
					}),
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
