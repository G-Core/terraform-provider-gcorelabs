package gcore

import (
	"fmt"
	"os"
	"regexp"
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

	fullName := "gcore_network.acctest"

	create := Params{
		Name: "create_test",
		Mtu:  1450,
		Type: "vxlan",
	}

	update_name := Params{
		Name: "update_test",
	}

	update_mtu := Params{
		Mtu: 1300,
	}

	update_type := Params{
		Type: "vlan",
	}

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
				Config:      NetworkTemplate(&update_mtu),
				ExpectError: regexp.MustCompile(`[Update a Network [0-9a-zA-Z\-]{36}] Validation error: unable to update 'mtu' field because it is immutable`),
			},
			{
				Config:      NetworkTemplate(&update_type),
				ExpectError: regexp.MustCompile(`[Update a Network [0-9a-zA-Z\-]{36}] Validation error: unable to update 'type' field because it is immutable`),
			},
		},
	})
}

func TestAccImportNetwork(t *testing.T) {
	fullName := "gcore_network.import_acctest"
	importStateIDPrefix := fmt.Sprintf("%s:%s:", os.Getenv("TEST_PROJECT_ID"), os.Getenv("TEST_REGION_ID"))

	var NetworkTemplate = fmt.Sprintf(`
		resource "gcore_network" "import_acctest" {
   			name = "import_test"
           	mtu = 1200
           	type = "vxlan"
			%s
			%s
       	}
		`, regionInfo(), projectInfo())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccNetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config: NetworkTemplate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(fullName),
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
	client, err := CreateTestClient(config.Provider, networksPoint)
	if err != nil {
		return err
	}
	allPages, err := networks.List(client).AllPages()
	if err != nil {
		return err
	}
	allNetworks, err := networks.ExtractNetworks(allPages)
	if err != nil {
		return err
	}
	if len(allNetworks) > 0 {
		return fmt.Errorf("Test client has networks: %v", allNetworks)
	}
	return nil
}
