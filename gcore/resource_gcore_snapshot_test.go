package gcore

import (
	"fmt"
	"os"
	"testing"

	"github.com/G-Core/gcorelabscloud-go/gcore/network/v1/networks"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccSnapshot(t *testing.T) {

	type Params struct {
		Name        string
		Description string
		Status      string
		Size        int
		VolumeID    string
	}

	create := Params{
		Name:     "test",
		VolumeID: GCORE_VOLUME_ID,
	}

	update := Params{
		Name:     "test",
		VolumeID: GCORE_VOLUME_ID,
	}

	fullName := "gcore_snapshot.acctest"
	importStateIDPrefix := fmt.Sprintf("%s:%s:", os.Getenv("TEST_PROJECT_ID"), os.Getenv("TEST_REGION_ID"))

	SnapshotTemplate := func(params *Params) string {

		additional := fmt.Sprintf("%s\n        %s", regionInfo(), projectInfo())

		template := fmt.Sprintf(`
		resource "gcore_snapshot" "acctest" {
			name = "%s"
			volume_id = "%s"
			%s
		`, params.Name, params.VolumeID, additional)

		return template + "\n}"
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckSnapshot(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccSnapshotDestroy,
		Steps: []resource.TestStep{
			{
				Config: SnapshotTemplate(&create),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(fullName),
					resource.TestCheckResourceAttr(fullName, "volume_id", create.VolumeID),
				),
			},
			{
				Config: SnapshotTemplate(&update),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(fullName),
					resource.TestCheckResourceAttr(fullName, "volume_id", update.VolumeID),
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

func testAccSnapshotDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	client, err := CreateTestClient(config.Provider, snapshotsPoint)
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gcore_snapshot" {
			continue
		}

		_, err := networks.Get(client, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Snapshot still exists")
		}
	}

	return nil
}
