package gcore

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"testing"

	"github.com/G-Core/gcorelabscloud-go/gcore/volume/v1/volumes"
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

	fullName := "gcore_volume.acctest"

	create := Params{
		Name:   "create_test",
		Size:   1,
		Type:   "standard",
		Source: "new-volume",
	}

	update := Params{
		Name:   "create_test",
		Size:   2,
		Type:   "ssd_hiiops",
		Source: "new-volume",
	}

	update_name := Params{
		Name:   "update_test",
		Size:   2,
		Type:   "ssd_hiiops",
		Source: "new-volume",
	}

	update_source := Params{
		Name:   "create_test",
		Size:   2,
		Type:   "ssd_hiiops",
		Source: "update-volume",
	}

	update_snapshot := Params{
		Name:       "create_test",
		Size:       2,
		Type:       "ssd_hiiops",
		Source:     "new-volume",
		SnapshotID: "4aceaf03-a5b2-47b8-aad9-8feb655557a8",
	}

	update_image := Params{
		Name:    "create_test",
		Size:    2,
		Type:    "ssd_hiiops",
		Source:  "new-volume",
		ImageID: "4aceaf03-a5b2-47b8-aad9-8feb655557a8",
	}

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
			source = "%s"
			type_name = "%s"
			%s
		`, params.Name, params.Size, params.Source, params.Type, additional)

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
				Config:      VolumeTemplate(&update_name),
				ExpectError: regexp.MustCompile(`\[Update a volume [0-9a-zA-Z\-]{36}] Validation error: unable to update 'name' field because it is immutable`),
			},
			{
				Config:      VolumeTemplate(&update_source),
				ExpectError: regexp.MustCompile(`\[Update a volume [0-9a-zA-Z\-]{36}] Validation error: unable to update 'source' field because it is immutable`),
			},
			{
				Config:      VolumeTemplate(&update_snapshot),
				ExpectError: regexp.MustCompile(`\[Update a volume [0-9a-zA-Z\-]{36}] Validation error: unable to update 'snapshot_id' field because it is immutable`),
			},
			{
				Config:      VolumeTemplate(&update_image),
				ExpectError: regexp.MustCompile(`[Update a volume [0-9a-zA-Z\-]{36}] Validation error: unable to update 'image_id' field because it is immutable`),
			},
		},
	})
}

func TestAccImportVolume(t *testing.T) {
	fullName := "gcore_volume.import_acctest"
	importStateIDPrefix := fmt.Sprintf("%s:%s:", os.Getenv("TEST_PROJECT_ID"), os.Getenv("TEST_REGION_ID"))

	var VolumeTemplate = fmt.Sprintf(`
		resource "gcore_volume" "import_acctest" {
			name = "import_test"
			size = 1
			source = "new-volume"
			type_name = "standard"
			%s
			%s
       	}
		`, regionInfo(), projectInfo())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: VolumeTemplate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(fullName),
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
	client, err := CreateTestClient(config.Provider, volumesPoint)
	if err != nil {
		return err
	}
	allPages, err := volumes.List(client, nil).AllPages()
	if err != nil {
		return err
	}
	allVolumes, err := volumes.ExtractVolumes(allPages)
	if err != nil {
		return err
	}
	if len(allVolumes) > 0 {
		return fmt.Errorf("Test client has volumes: %v", allVolumes)
	}
	return nil
}
