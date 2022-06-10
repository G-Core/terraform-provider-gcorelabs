//go:build cloud
// +build cloud

package gcore

import (
	"fmt"
	"testing"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/faas/v1/faas"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccFaaSNamespaceDataSource(t *testing.T) {
	cfg, err := createTestConfig()
	if err != nil {
		t.Fatal(err)
	}

	client, err := CreateTestClient(cfg.Provider, faasPoint, versionPointV1)
	if err != nil {
		t.Fatal(err)
	}

	nsName := "test-ns"
	if err := createTestNamespace(client, nsName); err != nil {
		t.Fatal(err)
	}
	defer faas.DeleteNamespace(client, nsName)

	ns, err := faas.GetNamespace(client, nsName).Extract()
	if err != nil {
		t.Fatal(err)
	}

	fullName := "data.gcore_faas_namespace.acctest"
	tpl := func(n string) string {
		return fmt.Sprintf(`
			data "gcore_faas_namespace" "acctest" {
			  %s
              %s
              name = "%s"
			}
		`, projectInfo(), regionInfo(), n)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: tpl(nsName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(fullName),
					resource.TestCheckResourceAttr(fullName, "id", ns.Name),
					resource.TestCheckResourceAttr(fullName, "name", ns.Name),
				),
			},
		},
	})
}

func createTestNamespace(client *gcorecloud.ServiceClient, nsName string) error {
	opts := faas.CreateNamespaceOpts{
		Name: nsName,
		Envs: map[string]string{},
	}

	res, err := faas.CreateNamespace(client, opts).Extract()
	if err != nil {
		return err
	}

	taskID := res.Tasks[0]
	_, err = tasks.WaitTaskAndReturnResult(client, taskID, true, FloatingIPCreateTimeout, func(task tasks.TaskID) (interface{}, error) {
		_, err := tasks.Get(client, string(task)).Extract()
		if err != nil {
			return nil, fmt.Errorf("cannot get task with ID: %s. Error: %w", task, err)
		}
		return nil, nil
	})
	return err
}
