package gcore

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccStorage(t *testing.T) {

	random := time.Now().Nanosecond()
	alias := fmt.Sprintf("terraform_test_alias_%d_s3", random)
	resourceName := fmt.Sprintf("gcore_storage.terraform_test_%d_s3", random)

	templateCreate := func() string {
		return fmt.Sprintf(`
resource "gcore_storage" "terraform_test_%d_s3" {
  name = "terraform_test_%d"
  location = "s-ed1"
  type = "s3"
}
		`, random, random)
	}

	templateUpdate := func() string {
		return fmt.Sprintf(`
resource "gcore_storage" "terraform_test_%d_s3" {
  name = "terraform_test_%d"
  location = "s-ed1"
  type = "s3"
  server_alias = "%s"
}
		`, random, random, alias)
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckVars(t, GCORE_USERNAME_VAR, GCORE_PASSWORD_VAR, GCORE_STORAGE_URL_VAR)
		},
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: templateCreate(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, StorageSchemaLocation, "s-ed1"),
					resource.TestCheckResourceAttr(resourceName, StorageSchemaType, "s3"),
				),
			},
			{
				Config: templateUpdate(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, StorageSchemaServerAlias, alias),
				),
			},
		},
	})
}
