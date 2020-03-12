package main

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccCreateVolumeV1(t *testing.T) {

	name:="foo"
	fullName := fmt.Sprintf("gcore_volume.%s", name)
	size := 1
	typeName := "standard"
	newSize := 2
	newTypeName := "ssd_hiiops"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVolumeTemplate(name, size, typeName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(fullName),
					resource.TestCheckResourceAttr(fullName, "size", strconv.Itoa(size)),
					resource.TestCheckResourceAttr(fullName, "type_name", typeName),
				),
			},
			{
				Config: testAccVolumeTemplate(name, newSize, typeName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(fullName),
					resource.TestCheckResourceAttr(fullName, "size", strconv.Itoa(newSize)),
				),
			},
			{
				Config: testAccVolumeTemplate(name, newSize, newTypeName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(fullName),
					resource.TestCheckResourceAttr(fullName, "type_name", newTypeName),
				),
			},
		},
	})
}

func testAccCheckResourceExists(resourceName string) resource.TestCheckFunc {
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

func testAccVolumeTemplate(name string, size int, typeName string) string {
    return fmt.Sprintf(`
	%s
	
	resource "gcore_volume" "%s" {
		name = "%s"
		size = %d
		type_name = "%s"
		%s
		%s
	}
	`, providerData(), name, name, size, typeName, regionInfo(), projectInfo())
}
