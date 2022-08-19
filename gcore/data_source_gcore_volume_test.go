//go:build cloud
// +build cloud

package gcore

import (
	"fmt"
	"strconv"
	"testing"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
	"github.com/G-Core/gcorelabscloud-go/gcore/volume/v1/volumes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const (
	volumeTestName  = "test-volume"
	volume1TestName = "test-volume-1"
	volume2TestName = "test-volume-2"
	volumeTestSize  = 1
)

func TestAccVolumeDataSource(t *testing.T) {
	cfg, err := createTestConfig()
	if err != nil {
		t.Fatal(err)
	}

	client, err := CreateTestClient(cfg.Provider, volumesPoint, versionPointV1)
	if err != nil {
		t.Fatal(err)
	}

	opts1 := volumes.CreateOpts{
		Name:     volume1TestName,
		Size:     volumeTestSize,
		Source:   volumes.NewVolume,
		TypeName: volumes.Standard,
		Metadata: map[string]string{"key1": "val1", "key2": "val2"},
	}

	volume1ID, err := createTestVolume(client, opts1)
	if err != nil {
		t.Fatal(err)
	}

	opts2 := volumes.CreateOpts{
		Name:     volume2TestName,
		Size:     volumeTestSize,
		Source:   volumes.NewVolume,
		TypeName: volumes.Standard,
		Metadata: map[string]string{"key1": "val1", "key3": "val3"},
	}

	volume2ID, err := createTestVolume(client, opts2)
	if err != nil {
		t.Fatal(err)
	}

	defer volumes.Delete(client, volume1ID, volumes.DeleteOpts{})
	defer volumes.Delete(client, volume2ID, volumes.DeleteOpts{})

	fullName := "data.gcore_volume.acctest"
	tpl := func(name string, metaQuery string) string {
		return fmt.Sprintf(`
			data "gcore_volume" "acctest" {
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
				Config: tpl(opts1.Name, `metadata_k="key1"`),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(fullName),
					resource.TestCheckResourceAttr(fullName, "name", opts1.Name),
					resource.TestCheckResourceAttr(fullName, "id", volume1ID),
					resource.TestCheckResourceAttr(fullName, "size", strconv.Itoa(opts1.Size)),
					testAccCheckMetadata(fullName, true, map[string]string{
						"key1": "val1", "key2": "val2"}),
				),
			},
			{
				Config: tpl(opts2.Name, `metadata_kv={key3 = "val3"}`),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(fullName),
					resource.TestCheckResourceAttr(fullName, "name", opts2.Name),
					resource.TestCheckResourceAttr(fullName, "id", volume2ID),
					resource.TestCheckResourceAttr(fullName, "size", strconv.Itoa(opts2.Size)),
					testAccCheckMetadata(fullName, true, map[string]string{
						"key3": "val3",
					}),
				),
			},
		},
	})
}

func createTestVolume(client *gcorecloud.ServiceClient, opts volumes.CreateOpts) (string, error) {
	res, err := volumes.Create(client, opts).Extract()
	if err != nil {
		return "", err
	}

	taskID := res.Tasks[0]
	volumeID, err := tasks.WaitTaskAndReturnResult(client, taskID, true, volumeCreatingTimeout, func(task tasks.TaskID) (interface{}, error) {
		taskInfo, err := tasks.Get(client, string(task)).Extract()
		if err != nil {
			return nil, fmt.Errorf("cannot get task with ID: %s. Error: %w", task, err)
		}
		volumeID, err := volumes.ExtractVolumeIDFromTask(taskInfo)
		if err != nil {
			return nil, fmt.Errorf("cannot retrieve volume ID from task info: %w", err)
		}
		return volumeID, nil
	},
	)

	if err != nil {
		return "", err
	}
	return volumeID.(string), nil
}
