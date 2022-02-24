package gcore

import (
	"fmt"
	"testing"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore/loadbalancer/v1/lbpools"
	"github.com/G-Core/gcorelabscloud-go/gcore/loadbalancer/v1/listeners"
	"github.com/G-Core/gcorelabscloud-go/gcore/loadbalancer/v1/loadbalancers"
	"github.com/G-Core/gcorelabscloud-go/gcore/loadbalancer/v1/types"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const (
	lbPoolTestName = "test-lb-pool"
)

func TestAccLBMember(t *testing.T) {
	cfg, err := createTestConfig()
	if err != nil {
		t.Fatal(err)
	}

	client, err := CreateTestClient(cfg.Provider, LoadBalancersPoint, versionPointV1)
	if err != nil {
		t.Fatal(err)
	}

	clientListener, err := CreateTestClient(cfg.Provider, LBListenersPoint, versionPointV1)
	if err != nil {
		t.Fatal(err)
	}

	clientPool, err := CreateTestClient(cfg.Provider, LBPoolsPoint, versionPointV1)
	if err != nil {
		t.Fatal(err)
	}

	opts := loadbalancers.CreateOpts{
		Name: lbTestName,
		Listeners: []loadbalancers.CreateListenerOpts{{
			Name:         lbListenerTestName,
			ProtocolPort: 80,
			Protocol:     types.ProtocolTypeHTTP,
		}},
	}

	lbID, err := createTestLoadBalancerWithListener(client, opts)
	if err != nil {
		t.Fatal(err)
	}
	defer loadbalancers.Delete(client, lbID)

	ls, err := listeners.ListAll(clientListener, listeners.ListOpts{LoadBalancerID: &lbID})
	if err != nil {
		t.Fatal(err)
	}
	listener := ls[0]

	poolsOpts := lbpools.CreateOpts{
		Name:            lbPoolTestName,
		Protocol:        types.ProtocolTypeHTTP,
		LBPoolAlgorithm: types.LoadBalancerAlgorithmRoundRobin,
		LoadBalancerID:  lbID,
		ListenerID:      listener.ID,
	}
	poolID, err := createTestLBPool(clientPool, poolsOpts)

	type Params struct {
		Address string
		Port    string
		Weight  string
	}

	create := Params{"10.10.2.15", "8080", "1"}

	update := Params{"10.10.2.16", "8081", "5"}

	fullName := "gcore_lbmember.acctest"

	tpl := func(params *Params) string {
		return fmt.Sprintf(`
            resource "gcore_lbmember" "acctest" {
			  %s
              %s
			  pool_id = "%s"
			  address = "%s"
			  protocol_port = %s
			  weight = %s
			}
		`, projectInfo(), regionInfo(), poolID, params.Address, params.Port, params.Weight)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccLBMemberDestroy,
		Steps: []resource.TestStep{
			{
				Config: tpl(&create),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(fullName),
					resource.TestCheckResourceAttr(fullName, "address", create.Address),
					resource.TestCheckResourceAttr(fullName, "protocol_port", create.Port),
					resource.TestCheckResourceAttr(fullName, "weight", create.Weight),
				),
			},
			{
				Config: tpl(&update),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(fullName),
					resource.TestCheckResourceAttr(fullName, "address", update.Address),
					resource.TestCheckResourceAttr(fullName, "protocol_port", update.Port),
					resource.TestCheckResourceAttr(fullName, "weight", update.Weight),
				),
			},
		},
	})
}

func testAccLBMemberDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	client, err := CreateTestClient(config.Provider, LBPoolsPoint, versionPointV1)
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gcore_lbmember" {
			continue
		}

		pl, err := lbpools.Get(client, GCORE_LBPOOL_ID).Extract()
		if err != nil {
			switch err.(type) {
			case gcorecloud.ErrDefault404:
				return nil
			default:
				return err
			}
		}

		for _, m := range pl.Members {
			if rs.Primary.ID == m.ID {
				return fmt.Errorf("LBMember still exists")
			}
		}
	}

	return nil
}

func createTestLBPool(client *gcorecloud.ServiceClient, opts lbpools.CreateOpts) (string, error) {
	res, err := lbpools.Create(client, opts).Extract()
	if err != nil {
		return "", err
	}

	taskID := res.Tasks[0]
	poolID, err := tasks.WaitTaskAndReturnResult(client, taskID, true, LBPoolsCreateTimeout, func(task tasks.TaskID) (interface{}, error) {
		taskInfo, err := tasks.Get(client, string(task)).Extract()
		if err != nil {
			return nil, fmt.Errorf("cannot get task with ID: %s. Error: %w", task, err)
		}
		poolID, err := lbpools.ExtractPoolIDFromTask(taskInfo)
		if err != nil {
			return nil, fmt.Errorf("cannot retrieve lb pool ID from task info: %w", err)
		}
		return poolID, nil
	})
	if err != nil {
		return "", err
	}
	return poolID.(string), nil
}
