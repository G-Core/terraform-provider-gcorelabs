//go:build !cloud
// +build !cloud

package gcore

import (
	"context"
	"fmt"
	"testing"
	"time"

	key2 "github.com/G-Core/gcore-storage-sdk-go/swagger/client/key"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestStorageSFTPKeyDataSource(t *testing.T) {
	random := time.Now().Nanosecond()
	name := fmt.Sprintf("terraformtestsftpkey%d", random)

	dataSourceName := fmt.Sprintf("data.gcore_storage_sftp_key.%s_data", name)
	key := `ssh-rsa AAAAB3NzaC1yc2EAAAABIwAAAQEAklOUpkDHrfHY17SbrmTIpNLTGK9Tjom/BWDSUGPl+nafzlHDTYW7hdI4yZ5ew18JH4JW9jbhUFrviQzM7xlELEVf4h9lFX5QVkbPppSwg0cda3Pbv7kOdJ/MTyBlWXFCR+HAo3FXRitBqxiX1nKhXpHAZsMciLq8V6RjsNAQwdsdMFvSlVK/7XAt3FaoJoAsncM1Q9x5+3V0Ww68/eIFmb1zuUFljQJKprrX88XypNDvjYNby6vw/Pb0rwert/EnmZ+AW4OZPnTPI89ZPmVMLuayrD2cE86Z/il8b+gw3r3+1nKatmIkjn2so1d01QraTlMqVSsbxNrRFi9wrf+M7Q== schacon@mylaptop.local`

	cfg, err := createTestConfig()
	if err != nil {
		t.Fatal("create conf", err)
	}
	opts := []func(params *key2.KeyCreateHTTPParams){
		func(params *key2.KeyCreateHTTPParams) {
			params.Context = context.Background()
			params.Body.Name = name
			params.Body.Key = key
		},
	}
	k, err := cfg.StorageClient.CreateKey(opts...)
	if err != nil {
		t.Fatal("create key", err)
	}

	templateRead := func() string {
		return fmt.Sprintf(`
data "gcore_storage_sftp_key" "%s_data" {
  name = "%s"
}
		`, name, name)
	}

	defer func() {
		opts := []func(params *key2.KeyDeleteHTTPParams){
			func(params *key2.KeyDeleteHTTPParams) {
				params.Context = context.Background()
				params.ID = k.ID
			},
		}
		err = cfg.StorageClient.DeleteKey(opts...)
		if err != nil {
			t.Fatal("delete key", err)
		}
	}()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckVars(t, GCORE_USERNAME_VAR, GCORE_PASSWORD_VAR, GCORE_STORAGE_URL_VAR)
		},
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: templateRead(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(dataSourceName),
					resource.TestCheckResourceAttr(dataSourceName, StorageKeySchemaName, name),
					resource.TestCheckResourceAttr(dataSourceName, StorageKeySchemaId, fmt.Sprint(k.ID)),
				),
			},
		},
	})
}
