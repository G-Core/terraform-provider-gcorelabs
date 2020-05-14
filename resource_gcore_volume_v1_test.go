package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"testing"

	"github.com/G-Core/gcorelabscloud-go/gcore/volume/v1/volumes"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccCreateVolumeV1_(t *testing.T) {
	name := "create_test"
	fullName := fmt.Sprintf("gcore_volumeV1.%s", name)
	size := 1
	typeName := "standard"
	newSize := 2
	newTypeName := "ssd_hiiops"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVolumeTemplate(name, size, typeName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(fullName),
					resource.TestCheckResourceAttr(fullName, "size", strconv.Itoa(size)),
					resource.TestCheckResourceAttr(fullName, "type_name", typeName),
				),
			},
			{
				Config: testAccVolumeTemplate(name, newSize, typeName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(fullName),
					resource.TestCheckResourceAttr(fullName, "size", strconv.Itoa(newSize)),
				),
			},
			{
				Config: testAccVolumeTemplate(name, newSize, newTypeName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(fullName),
					resource.TestCheckResourceAttr(fullName, "type_name", newTypeName),
				),
			},
		},
	})
}

func TestAccVolumeV1UpdateChecker_(t *testing.T) {
	name := "invalid_update"
	fullName := fmt.Sprintf("gcore_volumeV1.%s", name)
	size := 1
	typeName := "standard"
	source := "new-volume"

	newName := fmt.Sprintf("%s%d", name, 2)
	newSource := "snapshot"
	newSnapshotID := "4aceaf03-a5b2-47b8-aad9-8feb655557a8"
	newImageID := "4aceaf03-a5b2-47b8-aad9-8feb655557a8"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVolumeFullTemplate(name, name, size, source, typeName, "", ""),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(fullName),
					resource.TestCheckResourceAttr(fullName, "size", strconv.Itoa(size)),
					resource.TestCheckResourceAttr(fullName, "type_name", typeName),
				),
			},
			{
				Config:      testAccVolumeFullTemplate(name, newName, size, source, typeName, "", ""),
				ExpectError: regexp.MustCompile(`\[Update a volume [0-9a-zA-Z\-]{36}] Validation error: unable to update name field \(from invalid_update to invalid_update2\) because it is immutable`),
			},
			{
				Config:      testAccVolumeFullTemplate(name, name, size, newSource, typeName, "", ""),
				ExpectError: regexp.MustCompile(`\[Update a volume [0-9a-zA-Z\-]{36}] Validation error: unable to update source field \(from new-volume to snapshot\) because it is immutable`),
			},
			{
				Config:      testAccVolumeFullTemplate(name, name, size, source, typeName, newSnapshotID, ""),
				ExpectError: regexp.MustCompile(`\[Update a volume [0-9a-zA-Z\-]{36}] Validation error: unable to update snapshot_id field \(from  to 4aceaf03-a5b2-47b8-aad9-8feb655557a8\) because it is immutable`),
			},
			{
				Config:      testAccVolumeFullTemplate(name, name, size, source, typeName, "", newImageID),
				ExpectError: regexp.MustCompile(`[Update a volume [0-9a-zA-Z\-]{36}] Validation error: unable to update image_id field \(from  to 4aceaf03-a5b2-47b8-aad9-8feb655557a8\) because it is immutable`),
			},
		},
	})
}

func TestAccImportVolumeV1_(t *testing.T) {
	name := "import_test"
	fullName := fmt.Sprintf("gcore_volumeV1.%s", name)
	importStateIDPrefix := fmt.Sprintf("%s:%s:", os.Getenv("TEST_PROJECT_ID"), os.Getenv("TEST_REGION_ID"))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVolumeTemplate(name, 1, "standard"),
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

func testAccCheckVolumeDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	provider := config.Provider
	client, err := CreateTestClient(provider)
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

func testAccCheckResourceExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// retrieve the resource by name from state
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Widget ID is not set")
		}
		return nil
	}
}

func testAccVolumeTemplate(name string, size int, typeName string) string {
	return fmt.Sprintf(`
	resource "gcore_volumeV1" "%s" {
		name = "%s"
		size = %d
		source = "new-volume"
		type_name = "%s"
		%s
		%s
	}
	`, name, name, size, typeName, regionInfo(), projectInfo())
}

func testAccVolumeFullTemplate(terrafomName string, name string, size int, source string, typeName string, snapshotID string, imageID string) string {
	additionalInfo := fmt.Sprintf("%s\n        %s", regionInfo(), projectInfo())
	if snapshotID != "" {
		additionalInfo += fmt.Sprintf(`%s        snapshot_id = "%s"`, "\n", snapshotID)
	}
	if imageID != "" {
		additionalInfo += fmt.Sprintf(`%s        image_id = "%s"`, "\n", imageID)
	}

	return fmt.Sprintf(`
	resource "gcore_volumeV1" "%s" {
		name = "%s"
		size = %d
		source = "%s"
		type_name = "%s"
		%s
	}
	`, terrafomName, name, size, source, typeName, additionalInfo)
}
