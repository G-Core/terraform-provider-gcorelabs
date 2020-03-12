package main

import (
	"fmt"
	"testing"

	//"git.gcore.com/terraform-provider-gcore/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccVolume(t *testing.T) {

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckGcoreVolumeConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckExampleResourceExists("gcore_volume.foo"),
				),
			},
		},
	})
}

func testAccCheckExampleResourceExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// retrieve the resource by name from state
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Widget ID is not set")
		}
		return nil
	}
}

const testAccCheckGcoreVolumeConfig = `
provider "gcore" {
	jwt = "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJ0b2tlbl90eXBlIjoiYWNjZXNzIiwiZXhwIjoyNTgwOTM4Nzg4LCJqdGkiOiJhZWQ2ZjQwNjhhYzM0NWNkYWM1MTcwZjk0MzcwMDIzMyIsInVzZXJfaWQiOjEsInVzZXJfdHlwZSI6InN5c3RlbV9hZG1pbiIsInVzZXJfZ3JvdXBzIjpudWxsLCJjbGllbnRfaWQiOm51bGwsImVtYWlsIjoidGVzdEB0ZXN0LnRlc3QiLCJ1c2VybmFtZSI6InRlc3RAdGVzdC50ZXN0IiwiaXNfYWRtaW4iOnRydWUsImNsaWVudF9uYW1lIjoidGVzdCIsInJqdGkiOiJjZmJhNzMxODhlOTg0MzgxODAzZDdmYzU3OWJmZWIxYyJ9.0xjny_NM1uLQ5gRT8ZSmA_tvyeNZs8BrPjSFfhkKJbk"
}

resource "gcore_volume" "foo" {
	name = 156
	size = 2
	type_name = "ssd_hiiops"
	region_id = 1
	project_id = 78
  }
`
