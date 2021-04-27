package gcore

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccCDNRule(t *testing.T) {
	fullName := "gcore_cdn_rule.acctest"

	type Params struct {
		Name    string
		Pattern string
	}

	create := Params{"All images", "/folder/images/*.png"}
	update := Params{"All scripts", "/folder/scripts/*.js"}

	template := func(params *Params) string {
		return fmt.Sprintf(`
resource "gcore_cdn_rule" "acctest" {
  resource_id = %s
  name = "%s"
  rule = "%s"
  rule_type = 0
}
		`, GCORE_CDN_RESOURCE_ID, params.Name, params.Pattern)
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckVars(t, GCORE_USERNAME_VAR, GCORE_PASSWORD_VAR, GCORE_CDN_URL_VAR, GCORE_CDN_RESOURCE_ID_VAR)
		},
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: template(&create),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(fullName),
					resource.TestCheckResourceAttr(fullName, "name", create.Name),
					resource.TestCheckResourceAttr(fullName, "rule", create.Pattern),
				),
			},
			{
				Config: template(&update),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(fullName),
					resource.TestCheckResourceAttr(fullName, "name", update.Name),
					resource.TestCheckResourceAttr(fullName, "rule", update.Pattern),
				),
			},
		},
	})
}
