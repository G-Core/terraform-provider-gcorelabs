package gcore

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/G-Core/gcorelabs-storage-sdk-go/swagger/client/storage"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
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
		CheckDestroy: func(s *terraform.State) error {
			config := testAccProvider.Meta().(*Config)
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			for _, rs := range s.RootModule().Resources {
				if rs.Type != "gcore_storage" {
					continue
				}
				opts := []func(opt *storage.StorageListHTTPV2Params){
					func(opt *storage.StorageListHTTPV2Params) { opt.Context = ctx },
					func(opt *storage.StorageListHTTPV2Params) { opt.ID = &rs.Primary.ID },
				}
				storages, err := config.StorageClient.StoragesList(opts...)
				if err != nil {
					return fmt.Errorf("find storage: %w", err)
				}
				if len(storages) == 0 {
					return nil
				}
				if storages[0].ProvisioningStatus == "ok" {
					return fmt.Errorf("storage #%s wasn't deleted correctrly", rs.Primary.ID)
				}
			}
			return nil
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
