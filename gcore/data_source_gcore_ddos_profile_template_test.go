//go:build cloud
// +build cloud

package gcore

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/G-Core/gcorelabscloud-go/gcore/ddos/v1/ddos"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDDoSProtectionProfileTemplatesTest(t *testing.T) {
	cfg, err := createTestConfig()
	if err != nil {
		t.Fatal(err)
	}

	client, err := CreateTestClient(cfg.Provider, ddosTemplatesPoint, versionPointV1)
	if err != nil {
		t.Fatal(err)
	}

	templates, err := ddos.ListAllProfileTemplates(client)
	if err != nil {
		t.Fatal(err)
	}

	if len(templates) == 0 {
		t.Fatal("templates not found: templates list empty")
	}

	var template *ddos.ProfileTemplate
	for _, tmp := range templates {
		if len(tmp.Fields) > 0 {
			template = &tmp
			break
		}
	}

	if template == nil {
		t.Fatal("templates not found: there are no templates with non-empty fields")
	}

	fullName := "data.gcore_ddos_profile_template.acctest"
	tplByName := func(name string) string {
		return fmt.Sprintf(`
		data "gcore_ddos_profile_template" "acctest" {
			%s
			%s
			name = "%s"
		}
		`, projectInfo(), regionInfo(), name)
	}
	tplByID := func(id int) string {
		return fmt.Sprintf(`
		data "gcore_ddos_profile_template" "acctest" {
			%s
			%s
			template_id = "%d"
		}
		`, projectInfo(), regionInfo(), id)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: tplByID(template.ID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(fullName),
					resource.TestCheckResourceAttr(fullName, "name", template.Name),
					resource.TestCheckResourceAttr(fullName, "id", strconv.Itoa(template.ID)),
				),
			},
			{
				Config: tplByName(template.Name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(fullName),
					resource.TestCheckResourceAttr(fullName, "name", template.Name),
					resource.TestCheckResourceAttr(fullName, "id", strconv.Itoa(template.ID)),
				),
			},
		},
	})
}
