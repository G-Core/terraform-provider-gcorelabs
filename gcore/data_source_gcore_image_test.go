//go:build cloud
// +build cloud

package gcore

import (
	"fmt"
	"strconv"
	"testing"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/image/v1/images"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func uploadTestImage(client *gcorecloud.ServiceClient, opts images.UploadOpts) (string, error) {
	res, err := images.Upload(client, opts).Extract()
	if err != nil {
		return "", err
	}

	taskID := res.Tasks[0]
	imageID, err := tasks.WaitTaskAndReturnResult(client, taskID, true, ImageUploadTimeout, func(task tasks.TaskID) (interface{}, error) {
		taskInfo, err := tasks.Get(client, string(task)).Extract()
		if err != nil {
			return nil, fmt.Errorf("cannot get task with ID: %s. Error: %w", task, err)
		}
		id, err := images.ExtractImageIDFromTask(taskInfo)
		if err != nil {
			return nil, fmt.Errorf("cannot retrieve Image ID from task info: %w", err)
		}
		return id, nil
	})

	if err != nil {
		return "", err
	}
	return imageID.(string), nil
}

func TestAccImageDataSource(t *testing.T) {
	cfg, err := createTestConfig()
	if err != nil {
		t.Fatal(err)
	}

	client, err := CreateTestClient(cfg.Provider, imagesPoint, versionPointV1)
	if err != nil {
		t.Fatal(err)
	}
	downloadClient, err := CreateTestClient(cfg.Provider, downloadImagePoint, versionPointV1)
	if err != nil {
		t.Fatal(err)
	}

	opts1 := images.UploadOpts{
		HwMachineType:  "q35",
		SshKey:         "allow",
		Name:           "test_image_tf1",
		OSType:         "linux",
		URL:            "http://mirror.noris.net/cirros/0.4.0/cirros-0.4.0-x86_64-disk.img",
		HwFirmwareType: "uefi",
		Metadata:       map[string]string{"key1": "val1", "key2": "val2"},
	}

	opts2 := opts1
	opts2.Name = "test_image_tf2"
	opts2.Metadata = map[string]string{"key1": "val1", "key3": "val3"}

	image1ID, err := uploadTestImage(downloadClient, opts1)
	if err != nil {
		t.Fatal(err)
	}
	defer images.Delete(client, image1ID)

	image2ID, err := uploadTestImage(downloadClient, opts2)
	if err != nil {
		t.Fatal(err)
	}
	defer images.Delete(client, image2ID)

	image1, err := images.Get(client, image1ID).Extract()
	if err != nil {
		t.Fatal(err)
	}

	image2, err := images.Get(client, image2ID).Extract()
	if err != nil {
		t.Fatal(err)
	}

	fullName := "data.gcore_image.acctest"
	tpl := func(name string, metaQuery string) string {
		return fmt.Sprintf(`
			data "gcore_image" "acctest" {
			  %s
              %s
              name = "%s"
   				%s
			}
		`, projectInfo(), regionInfo(), name, metaQuery)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: tpl(image1.Name, `metadata_k="key1"`),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(fullName),
					resource.TestCheckResourceAttr(fullName, "name", image1.Name),
					resource.TestCheckResourceAttr(fullName, "id", image1.ID),
					resource.TestCheckResourceAttr(fullName, "min_disk", strconv.Itoa(image1.MinDisk)),
					resource.TestCheckResourceAttr(fullName, "min_ram", strconv.Itoa(image1.MinRAM)),
					resource.TestCheckResourceAttr(fullName, "os_distro", image1.OsDistro),
					resource.TestCheckResourceAttr(fullName, "os_version", image1.OsVersion),
					testAccCheckMetadata(fullName, true, map[string]string{
						"key1": "val1", "key2": "val2"}),
				),
			},
			{
				Config: tpl(image2.Name, `metadata_kv={key3 = "val3"}`),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(fullName),
					resource.TestCheckResourceAttr(fullName, "name", image2.Name),
					resource.TestCheckResourceAttr(fullName, "id", image2.ID),
					resource.TestCheckResourceAttr(fullName, "min_disk", strconv.Itoa(image2.MinDisk)),
					resource.TestCheckResourceAttr(fullName, "min_ram", strconv.Itoa(image2.MinRAM)),
					resource.TestCheckResourceAttr(fullName, "os_distro", image2.OsDistro),
					resource.TestCheckResourceAttr(fullName, "os_version", image2.OsVersion),
					testAccCheckMetadata(fullName, true, map[string]string{
						"key3": "val3",
					}),
				),
			},
		},
	})
}
