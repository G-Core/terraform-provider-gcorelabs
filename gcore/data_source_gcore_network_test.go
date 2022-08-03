//go:build cloud
// +build cloud

package gcore

import (
	"fmt"
	"testing"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/network/v1/networks"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const (
	network1TestName = "test-network1"
	network2TestName = "test-network2"
	// Used in other modules
	networkTestName = "test-network"
)

func TestAccNetworkDataSource(t *testing.T) {
	cfg, err := createTestConfig()
	if err != nil {
		t.Fatal(err)
	}

	client, err := CreateTestClient(cfg.Provider, networksPoint, versionPointV1)
	if err != nil {
		t.Fatal(err)
	}

	opts1 := networks.CreateOpts{
		Name:     network1TestName,
		Metadata: map[string]string{"key1": "val1", "key2": "val2"},
	}

	network1ID, err := createTestNetwork(client, opts1)
	if err != nil {
		t.Fatal(err)
	}
	opts2 := networks.CreateOpts{
		Name:     network2TestName,
		Metadata: map[string]string{"key1": "val1", "key3": "val3"},
	}

	network2ID, err := createTestNetwork(client, opts2)
	if err != nil {
		t.Fatal(err)
	}

	defer deleteTestNetwork(client, network1ID)
	defer deleteTestNetwork(client, network2ID)

	fullName := "data.gcore_network.acctest"
	tpl1 := func(name string) string {
		return fmt.Sprintf(`
			data "gcore_network" "acctest" {
			  %s
              %s
              name = "%s"
              metadata_k="key1"
			}
		`, projectInfo(), regionInfo(), name)
	}
	tpl2 := func(name string) string {
		return fmt.Sprintf(`
			data "gcore_network" "acctest" {
			  %s
              %s
              name = "%s"
 			  metadata_kv={
                  key3 = "val3"
			  }
			}
		`, projectInfo(), regionInfo(), name)
	}
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: tpl1(opts1.Name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(fullName),
					resource.TestCheckResourceAttr(fullName, "name", opts1.Name),
					resource.TestCheckResourceAttr(fullName, "id", network1ID),
					testAccCheckMetadata(fullName, true, map[string]string{
						"key1": "val1", "key2": "val2",
					}),
				),
			},
			{
				Config: tpl2(opts2.Name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(fullName),
					resource.TestCheckResourceAttr(fullName, "name", opts2.Name),
					resource.TestCheckResourceAttr(fullName, "id", network2ID),
					testAccCheckMetadata(fullName, true, map[string]string{
						"key3": "val3",
					}),
				),
			},
		},
	})
}

func createTestNetwork(client *gcorecloud.ServiceClient, opts networks.CreateOpts) (string, error) {
	res, err := networks.Create(client, opts).Extract()
	if err != nil {
		return "", err
	}

	taskID := res.Tasks[0]
	networkID, err := tasks.WaitTaskAndReturnResult(client, taskID, true, networkCreatingTimeout, func(task tasks.TaskID) (interface{}, error) {
		taskInfo, err := tasks.Get(client, string(task)).Extract()
		if err != nil {
			return nil, fmt.Errorf("cannot get task with ID: %s. Error: %w", task, err)
		}
		networkID, err := networks.ExtractNetworkIDFromTask(taskInfo)
		if err != nil {
			return nil, fmt.Errorf("cannot retrieve network ID from task info: %w", err)
		}
		return networkID, nil
	},
	)
	if err != nil {
		return "", err
	}
	return networkID.(string), nil
}

func deleteTestNetwork(client *gcorecloud.ServiceClient, networkID string) error {
	results, err := networks.Delete(client, networkID).Extract()
	if err != nil {
		return err
	}
	taskID := results.Tasks[0]
	_, err = tasks.WaitTaskAndReturnResult(client, taskID, true, networkDeleting, func(task tasks.TaskID) (interface{}, error) {
		_, err := networks.Get(client, networkID).Extract()
		if err == nil {
			return nil, fmt.Errorf("cannot delete network with ID: %s", networkID)
		}
		switch err.(type) {
		case gcorecloud.ErrDefault404:
			return nil, nil
		default:
			return nil, err
		}
	})
	return err
}
