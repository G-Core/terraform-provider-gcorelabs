//go:build !cloud
// +build !cloud

package gcore

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/G-Core/gcore-storage-sdk-go/swagger/client/storage"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccStorageS3Bucket(t *testing.T) {

	random := time.Now().Nanosecond()
	storageResourceName := fmt.Sprintf("gcore_storage_s3.terraform_test_%d_s3", random)
	bucketResourceName := fmt.Sprintf("gcore_storage_s3_bucket.terraform_test_%d_s3_bucket", random)
	name := fmt.Sprintf("terraform_test_%d", random)

	templateCreateBucket := func() string {
		return fmt.Sprintf(`
resource "gcore_storage_s3" "terraform_test_%d_s3" {
  name = "terraform_test_%d"
  location = "s-ed1"
}

resource "gcore_storage_s3_bucket" "terraform_test_%d_s3_bucket" {
  name = "terraform_test_%d"
  storage_id = %s.id
}
		`, random, random, random, random, storageResourceName)
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
				if rs.Type != "gcore_storage_s3" {
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
				Config: templateCreateBucket(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(bucketResourceName),
					resource.TestCheckResourceAttr(bucketResourceName, StorageS3BucketSchemaName, name),
				),
			},
		},
	})
}
