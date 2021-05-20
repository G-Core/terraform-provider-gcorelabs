package gcore

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccStorageKey(t *testing.T) {

	random := time.Now().Nanosecond()
	name := fmt.Sprintf("terraform_test_key_%d", random)
	resourceName := fmt.Sprintf("gcore_storage_key.terraform_test_%d_s3", random)
	key := `ssh-rsa AAAAB3NzaC1yc2EAAAABIwAAAQEAklOUpkDHrfHY17SbrmTIpNLTGK9Tjom/BWDSU
GPl+nafzlHDTYW7hdI4yZ5ew18JH4JW9jbhUFrviQzM7xlELEVf4h9lFX5QVkbPppSwg0cda3
Pbv7kOdJ/MTyBlWXFCR+HAo3FXRitBqxiX1nKhXpHAZsMciLq8V6RjsNAQwdsdMFvSlVK/7XA
t3FaoJoAsncM1Q9x5+3V0Ww68/eIFmb1zuUFljQJKprrX88XypNDvjYNby6vw/Pb0rwert/En
mZ+AW4OZPnTPI89ZPmVMLuayrD2cE86Z/il8b+gw3r3+1nKatmIkjn2so1d01QraTlMqVSsbx
NrRFi9wrf+M7Q== schacon@mylaptop.local`

	templateCreate := func() string {
		return fmt.Sprintf(`
resource "gcore_storage_key" "terraform_test_key_%d" {
  name = "%s"
  key = "%s"
}
		`, random, name, key)
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
					resource.TestCheckResourceAttr(resourceName, StorageKeySchemaName, name),
				),
			},
		},
	})
}
