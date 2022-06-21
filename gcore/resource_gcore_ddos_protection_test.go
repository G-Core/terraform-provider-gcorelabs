//go:build cloud
// +build cloud

package gcore

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/G-Core/gcorelabscloud-go/gcore/ddos/v1/ddos"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDDoSProfile(t *testing.T) {
	type Params struct {
		ProfileTemplate string
		ProfileFields   string
		BGP             bool
	}

	profileTmpl := func(params *Params) string {
		template := `
	resource "gcore_baremetal" "bm" {
		%[1]s
		%[2]s
		name = "baremetal_acctest"
		flavor_id = "bm1-hf-medium-fake"
		image_id = "570fb9a3-5074-4539-b0d0-ec49f8c463aa"
		interface {
			type = "external"
			is_parent = "true"
		}
	}

	resource "gcore_ddos_protection" "acctest" {
		%[1]s
		%[2]s
		ip_address = gcore_baremetal.bm.addresses.0.net.0.addr
		profile_template = %[3]s
		bm_instance_id = gcore_baremetal.bm.id
		active = true
		bgp = %[4]t
		fields {
			base_field  = 118
			field_value = "%[5]s"
		}
	}
`
		return fmt.Sprintf(template, projectInfo(), regionInfo(), params.ProfileTemplate, params.BGP, params.ProfileFields)
	}

	createParams := Params{
		ProfileTemplate: "63",
		ProfileFields:   "[33033]",
		BGP:             false,
	}

	updateParams := Params{
		ProfileTemplate: createParams.ProfileTemplate,
		ProfileFields:   "[33031,33034]",
		BGP:             true,
	}

	fullName := "gcore_ddos_protection.acctest"
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      checkDestroyDDoSProtectionProfile,
		Steps: []resource.TestStep{
			{
				Config: profileTmpl(&createParams),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(fullName),
					resource.TestCheckResourceAttr(fullName, "profile_template", createParams.ProfileTemplate),
					resource.TestCheckResourceAttrSet(fullName, "ip_address"),
					resource.TestCheckResourceAttr(fullName, "active", "true"),
					resource.TestCheckResourceAttr(fullName, "bgp", "false"),
					resource.TestCheckResourceAttr(fullName, "fields.0.field_value", createParams.ProfileFields),
				),
			},
			{
				Config: profileTmpl(&updateParams),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(fullName),
					resource.TestCheckResourceAttr(fullName, "profile_template", updateParams.ProfileTemplate),
					resource.TestCheckResourceAttrSet(fullName, "ip_address"),
					resource.TestCheckResourceAttr(fullName, "active", "true"),
					resource.TestCheckResourceAttr(fullName, "bgp", "true"),
					resource.TestCheckResourceAttr(fullName, "fields.0.field_value", updateParams.ProfileFields),
				),
			},
		},
	})
}

func checkDestroyDDoSProtectionProfile(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	client, err := CreateTestClient(config.Provider, ddosProfilePoint, versionPointV1)
	if err != nil {
		return err
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gcore_ddos_protection" {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return err
		}

		profiles, err := ddos.ListAllProfiles(client)
		if err != nil {
			return err
		}

		for _, profile := range profiles {
			if profile.ID == id {
				return fmt.Errorf("ddos protection profile still exists")
			}
		}
	}

	return nil
}
