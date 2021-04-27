package gcore

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOriginGroup(t *testing.T) {
	fullName := "gcore_cdn_origingroup.acctest"

	type Params struct {
		Source  string
		Enabled string
	}

	create := Params{"google.com", "true"}
	update := Params{"tut.by", "false"}

	template := func(params *Params) string {
		return fmt.Sprintf(`
            resource "gcore_cdn_origingroup" "acctest" {
			  name = terraform_acctest_group
			  use_next = true

			  origin {
			    source = "%s"
				enabled = %s
			  }

			  origin {
			    source = "yandex.ru"
			    enabled = true
			    backup = true
			  }
			}
		`, params.Source, params.Enabled)
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckVars(t, GCORE_USERNAME_VAR, GCORE_PASSWORD_VAR, GCORE_CDN_URL_VAR)
		},
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: template(&create),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(fullName),
					resource.TestCheckResourceAttr(fullName, "name", "terraform_acctest_group"),
					resource.TestCheckResourceAttr(fullName, "source", create.Source),
					resource.TestCheckResourceAttr(fullName, "enabled", create.Enabled),
				),
			},
			{
				Config: template(&update),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(fullName),
					resource.TestCheckResourceAttr(fullName, "name", "terraform_acctest_group"),
					resource.TestCheckResourceAttr(fullName, "source", update.Source),
					resource.TestCheckResourceAttr(fullName, "enabled", update.Enabled),
				),
			},
		},
	})
}
