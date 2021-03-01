package gcore

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/G-Core/gcorelabscloud-go/gcore/image/v1/images"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccImageDataSource(t *testing.T) {
	cfg, err := createTestConfig()
	if err != nil {
		t.Fatal(err)
	}

	client, err := CreateTestClient(cfg.Provider, imagesPoint, versionPointV1)
	if err != nil {
		t.Fatal(err)
	}

	imgs, err := images.ListAll(client, images.ListOpts{})
	if err != nil {
		t.Fatal(err)
	}

	if len(imgs) == 0 {
		t.Fatal("images not found")
	}

	img := imgs[0]

	fullName := "data.gcore_image.acctest"
	tpl := func(name string) string {
		return fmt.Sprintf(`
			data "gcore_image" "acctest" {
			  %s
              %s
              name = "%s"
			}
		`, projectInfo(), regionInfo(), name)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: tpl(img.Name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(fullName),
					resource.TestCheckResourceAttr(fullName, "name", img.Name),
					resource.TestCheckResourceAttr(fullName, "id", img.ID),
					resource.TestCheckResourceAttr(fullName, "min_disk", strconv.Itoa(img.MinDisk)),
					resource.TestCheckResourceAttr(fullName, "min_ram", strconv.Itoa(img.MinRAM)),
					resource.TestCheckResourceAttr(fullName, "os_distro", img.OsDistro),
					resource.TestCheckResourceAttr(fullName, "os_version", img.OsVersion),
				),
			},
		},
	})
}
