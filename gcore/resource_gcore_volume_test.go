package gcore

import (
	"fmt"
	"github.com/G-Core/gcorelabscloud-go/gcore/network/v1/networks"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccVolume(t *testing.T) {

	type Params struct {
		Name       string
		Size       int
		Type       string
		Source     string
		SnapshotID string
		ImageID    string
	}

	create := Params{
		Name: "test",
		Size: 1,
		Type: "standard",
	}

	update := Params{
		Name: "test",
		Size: 2,
		Type: "ssd_hiiops",
	}

	fullName := "gcore_volume.acctest"
	importStateIDPrefix := fmt.Sprintf("%s:%s:", os.Getenv("TEST_PROJECT_ID"), os.Getenv("TEST_REGION_ID"))

	VolumeTemplate := func(params *Params) string {

		additional := fmt.Sprintf("%s\n        %s", regionInfo(), projectInfo())
		if params.SnapshotID != "" {
			additional += fmt.Sprintf(`%s        snapshot_id = "%s"`, "\n", params.SnapshotID)
		}
		if params.ImageID != "" {
			additional += fmt.Sprintf(`%s        image_id = "%s"`, "\n", params.ImageID)
		}

		template := fmt.Sprintf(`
		resource "gcore_volume" "acctest" {
			name = "%s"
			size = %d
			type_name = "%s"
			%s
		`, params.Name, params.Size, params.Type, additional)

		return template + "\n}"
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: VolumeTemplate(&create),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(fullName),
					resource.TestCheckResourceAttr(fullName, "size", strconv.Itoa(create.Size)),
					resource.TestCheckResourceAttr(fullName, "type_name", create.Type),
				),
			},
			{
				Config: VolumeTemplate(&update),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(fullName),
					resource.TestCheckResourceAttr(fullName, "size", strconv.Itoa(update.Size)),
					resource.TestCheckResourceAttr(fullName, "type_name", update.Type),
				),
			},
			{
				ImportStateIdPrefix: importStateIDPrefix,
				ResourceName:        fullName,
				ImportState:         true,
			},
		},
	})
}

func testAccVolumeDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	client, err := CreateTestClient(config.Provider, volumesPoint, versionPointV1)
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gcore_volume" {
			continue
		}

		_, err := networks.Get(client, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Volume still exists")
		}
	}

	return nil
}
