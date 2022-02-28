//go:build cloud
// +build cloud

package gcore

import (
	"fmt"
	"testing"

	"github.com/G-Core/gcorelabscloud-go/gcore/keypair/v2/keypairs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const pkTest = `ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC1bdbQYquD/swsZpFPXagY9KvhlNUTKYMdhRNtlGglAMgRxJS3Q0V74BNElJtP+UU/AbZD4H2ZAwW3PLLD/maclnLlrA48xg/ez9IhppBop0WADZ/nB4EcvQfR/Db7nHDTZERW6EiiGhV6CkHVasK2sY/WNRXqPveeWUlwCqtSnU90l/s9kQCoEfkM2auO6ppJkVrXbs26vcRclS8KL7Cff4HwdVpV7b+edT5seZdtrFUCbkEof9D9nGpahNvg8mYWf0ofx4ona4kaXm1NdPID+ljvE/dbYUX8WZRmyLjMvVQS+VxDJtsiDQIVtwbC4w+recqwDvHhLWwoeczsbEsp ondi@ds`

func TestAccKeyPair(t *testing.T) {
	type Params struct {
		Name string
		PK   string
	}

	create := Params{
		Name: "test",
		PK:   pkTest,
	}

	fullName := "gcore_keypair.acctest"

	kpTemplate := func(params *Params) string {
		return fmt.Sprintf(`
			resource "gcore_keypair" "acctest" {
			  %s
			  public_key = "%s"
			  sshkey_name = "%s"
			}
		`, projectInfo(), params.PK, params.Name)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccKeypairDestroy,
		Steps: []resource.TestStep{
			{
				Config: kpTemplate(&create),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(fullName),
					resource.TestCheckResourceAttr(fullName, "sshkey_name", create.Name),
					resource.TestCheckResourceAttr(fullName, "public_key", create.PK),
				),
			},
		},
	})
}

func testAccKeypairDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	client, err := CreateTestClient(config.Provider, keypairsPoint, versionPointV1)
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gcore_keypair" {
			continue
		}

		_, err := keypairs.Get(client, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("KeyPair still exists")
		}
	}

	return nil
}
