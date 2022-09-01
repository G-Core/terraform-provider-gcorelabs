//go:build cloud
// +build cloud

package gcore

import (
	"fmt"
	"testing"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/loadbalancer/v1/loadbalancers"
	"github.com/G-Core/gcorelabscloud-go/gcore/loadbalancer/v1/types"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const (
	lbTestName         = "test-lb"
	lbListenerTestName = "test-listener"

	lb1TestName = "test-lb1"
	lb2TestName = "test-lb2"

	lb1ListenerTestName = "test-listener1"
	lb2ListenerTestName = "test-listener2"
)

func TestAccLoadBalancerDataSource(t *testing.T) {
	cfg, err := createTestConfig()
	if err != nil {
		t.Fatal(err)
	}

	client, err := CreateTestClient(cfg.Provider, LoadBalancersPoint, versionPointV1)
	if err != nil {
		t.Fatal(err)
	}

	opts1 := loadbalancers.CreateOpts{
		Name: lb1TestName,
		Listeners: []loadbalancers.CreateListenerOpts{{
			Name:         lb1ListenerTestName,
			ProtocolPort: 80,
			Protocol:     types.ProtocolTypeHTTP,
		}},
		Metadata: map[string]string{"key1": "val1", "key2": "val2"},
	}

	opts2 := loadbalancers.CreateOpts{
		Name: lb2TestName,
		Listeners: []loadbalancers.CreateListenerOpts{{
			Name:         lb2ListenerTestName,
			ProtocolPort: 80,
			Protocol:     types.ProtocolTypeHTTP,
		}},
		Metadata: map[string]string{"key1": "val1", "key3": "val3"},
	}

	lb1ID, err := createTestLoadBalancerWithListener(client, opts1)
	if err != nil {
		t.Fatal(err)
	}

	lb2ID, err := createTestLoadBalancerWithListener(client, opts2)
	if err != nil {
		t.Fatal(err)
	}

	defer loadbalancers.Delete(client, lb1ID)
	defer loadbalancers.Delete(client, lb2ID)

	fullName := "data.gcore_loadbalancer.acctest"
	tpl := func(name string, metaQuery string) string {
		return fmt.Sprintf(`
			data "gcore_loadbalancer" "acctest" {
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
					resource.TestCheckResourceAttr(fullName, "id", lb1ID),
					testAccCheckMetadata(fullName, true, map[string]string{
						"key1": "val1", "key2": "val2"}),
				),
			},
			{
				Config: tpl(opts2.Name, `metadata_kv={key3 = "val3"}`),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(fullName),
					resource.TestCheckResourceAttr(fullName, "name", opts2.Name),
					resource.TestCheckResourceAttr(fullName, "id", lb2ID),
					testAccCheckMetadata(fullName, true, map[string]string{
						"key3": "val3"}),
				),
			},
		},
	})
}

func createTestLoadBalancerWithListener(client *gcorecloud.ServiceClient, opts loadbalancers.CreateOpts) (string, error) {
	res, err := loadbalancers.Create(client, opts).Extract()
	if err != nil {
		return "", err
	}

	taskID := res.Tasks[0]
	lbID, err := tasks.WaitTaskAndReturnResult(client, taskID, true, LoadBalancerCreateTimeout, func(task tasks.TaskID) (interface{}, error) {
		taskInfo, err := tasks.Get(client, string(task)).Extract()
		if err != nil {
			return nil, fmt.Errorf("cannot get task with ID: %s. Error: %w", task, err)
		}
		lbID, err := loadbalancers.ExtractLoadBalancerIDFromTask(taskInfo)
		if err != nil {
			return nil, fmt.Errorf("cannot retrieve LoadBalancer ID from task info: %w", err)
		}
		return lbID, nil
	})
	if err != nil {
		return "", err
	}
	return lbID.(string), nil
}
