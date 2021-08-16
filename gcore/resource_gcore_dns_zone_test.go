package gcore

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDnsZone1(t *testing.T) {

	random := time.Now().Nanosecond()
	name := fmt.Sprintf("terraformtestkey%d", random)
	zone := name + ".com"
	resourceName := fmt.Sprintf("%s.%s", DNSZoneResource, name)

	templateCreate := func() string {
		return fmt.Sprintf(`
resource "%s" "%s" {
  name = "%s"
}
		`, DNSZoneResource, name, zone)
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckVars(t, GCORE_USERNAME_VAR, GCORE_PASSWORD_VAR, GCORE_DNS_URL_VAR)
		},
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: templateCreate(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, DNSZoneSchemaName, zone),
				),
			},
		},
	})
}
