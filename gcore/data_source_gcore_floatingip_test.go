//go:build cloud
// +build cloud

package gcore

import (
	"fmt"
	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/floatingip/v1/floatingips"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccFloatingIPDataSource(t *testing.T) {
	cfg, err := createTestConfig()
	if err != nil {
		t.Fatal(err)
	}

	client, err := CreateTestClient(cfg.Provider, floatingIPsPoint, versionPointV1)
	if err != nil {
		t.Fatal(err)
	}

	opts1 := floatingips.CreateOpts{
		Metadata: map[string]string{"key1": "val1", "key2": "val2"},
	}

	floatingIPID1, err := createTestFloatingIP(client, opts1)
	if err != nil {
		t.Fatal(err)
	}

	opts2 := floatingips.CreateOpts{
		Metadata: map[string]string{"key1": "val1", "key3": "val3"},
	}

	floatingIPID2, err := createTestFloatingIP(client, opts2)
	if err != nil {
		t.Fatal(err)
	}

	defer floatingips.Delete(client, floatingIPID1)
	defer floatingips.Delete(client, floatingIPID2)

	fip1, err := floatingips.Get(client, floatingIPID1).Extract()
	if err != nil {
		t.Fatal(err)
	}

	fip2, err := floatingips.Get(client, floatingIPID2).Extract()
	if err != nil {
		t.Fatal(err)
	}

	fullName := "data.gcore_floatingip.acctest"
	tpl := func(ip string, metaQuery string) string {
		return fmt.Sprintf(`
			data "gcore_floatingip" "acctest" {
			  %s
              %s
              floating_ip_address = "%s"
			  %s
			}
		`, projectInfo(), regionInfo(), ip, metaQuery)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: tpl(fip1.FloatingIPAddress.String(), `metadata_k="key1"`),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(fullName),
					resource.TestCheckResourceAttr(fullName, "id", floatingIPID1),
					resource.TestCheckResourceAttr(fullName, "floating_ip_address", fip1.FloatingIPAddress.String()),
					testAccCheckMetadata(fullName, true, map[string]string{
						"key1": "val1", "key2": "val2"}),
				),
			},
			{
				Config: tpl(fip2.FloatingIPAddress.String(), `metadata_kv={key3 = "val3"}`),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(fullName),
					resource.TestCheckResourceAttr(fullName, "id", floatingIPID2),
					resource.TestCheckResourceAttr(fullName, "floating_ip_address", fip2.FloatingIPAddress.String()),
					testAccCheckMetadata(fullName, true, map[string]string{
						"key3": "val3",
					}),
				),
			},
		},
	})
}

func createTestFloatingIP(client *gcorecloud.ServiceClient, opts floatingips.CreateOpts) (string, error) {
	res, err := floatingips.Create(client, opts).Extract()
	if err != nil {
		return "", err
	}

	taskID := res.Tasks[0]
	floatingIPID, err := tasks.WaitTaskAndReturnResult(client, taskID, true, FloatingIPCreateTimeout, func(task tasks.TaskID) (interface{}, error) {
		taskInfo, err := tasks.Get(client, string(task)).Extract()
		if err != nil {
			return nil, fmt.Errorf("cannot get task with ID: %s. Error: %w", task, err)
		}
		floatingIPID, err := floatingips.ExtractFloatingIPIDFromTask(taskInfo)
		if err != nil {
			return nil, fmt.Errorf("cannot retrieve FloatingIP ID from task info: %w", err)
		}
		return floatingIPID, nil
	})

	if err != nil {
		return "", err
	}
	return floatingIPID.(string), nil
}
