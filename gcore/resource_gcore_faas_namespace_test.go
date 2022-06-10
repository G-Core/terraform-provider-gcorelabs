//go:build cloud
// +build cloud

package gcore

import (
	"fmt"
	"testing"

	"github.com/G-Core/gcorelabscloud-go/gcore/faas/v1/faas"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccFaaSNamespace(t *testing.T) {

	type Params struct {
		Name        string
		Description string
	}

	create := Params{
		Name:        "ns-name",
		Description: "description",
	}

	update := Params{
		Name:        "ns-name",
		Description: "changed description",
	}

	fullName := "gcore_faas_namespace.acctest"

	tpl := func(params *Params) string {
		template := fmt.Sprintf(`
		resource "gcore_faas_namespace" "acctest" {
			name = "%s"
			description = "%s"
			%s
			%s
		}`, params.Name, params.Description, regionInfo(), projectInfo())

		return template
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      CheckDestroyFaaSNamespace,
		Steps: []resource.TestStep{
			{
				Config: tpl(&create),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(fullName),
					resource.TestCheckResourceAttr(fullName, "name", create.Name),
					resource.TestCheckResourceAttr(fullName, "description", create.Description),
				),
			},
			{
				Config: tpl(&update),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(fullName),
					resource.TestCheckResourceAttr(fullName, "name", update.Name),
					resource.TestCheckResourceAttr(fullName, "description", update.Description),
				),
			},
		},
	})
}

func CheckDestroyFaaSNamespace(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	client, err := CreateTestClient(config.Provider, faasPoint, versionPointV1)
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gcore_faas_namespace" {
			continue
		}

		_, err := faas.GetNamespace(client, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("faas namespace still exists")
		}
	}

	return nil
}
