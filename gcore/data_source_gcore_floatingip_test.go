package gcore

import (
	"fmt"
	"testing"

	"github.com/G-Core/gcorelabscloud-go/gcore/floatingip/v1/floatingips"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
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

	opts := floatingips.CreateOpts{}

	res, err := floatingips.Create(client, opts).Extract()
	if err != nil {
		t.Fatal(err)
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
		t.Fatal(err)
	}

	defer floatingips.Delete(client, floatingIPID.(string))

	fip, err := floatingips.Get(client, floatingIPID.(string)).Extract()
	if err != nil {
		t.Fatal(err)
	}

	fullName := "data.gcore_floatingip.acctest"
	tpl := func(ip string) string {
		return fmt.Sprintf(`
			data "gcore_floatingip" "acctest" {
			  %s
              %s
              floating_ip_address = "%s"
			}
		`, projectInfo(), regionInfo(), ip)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: tpl(fip.FloatingIPAddress.String()),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(fullName),
					resource.TestCheckResourceAttr(fullName, "id", floatingIPID.(string)),
					resource.TestCheckResourceAttr(fullName, "floating_ip_address", fip.FloatingIPAddress.String()),
				),
			},
		},
	})
}
