//go:build cloud
// +build cloud

package gcore

import (
	"fmt"
	"testing"

	"github.com/G-Core/gcorelabscloud-go/gcore/loadbalancer/v1/lbpools"
	"github.com/G-Core/gcorelabscloud-go/gcore/loadbalancer/v1/listeners"
	"github.com/G-Core/gcorelabscloud-go/gcore/loadbalancer/v1/loadbalancers"
	"github.com/G-Core/gcorelabscloud-go/gcore/loadbalancer/v1/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccLBPool(t *testing.T) {
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

	type Params struct {
		Name        string
		LBAlgorithm string
	}

	create := Params{"test", "ROUND_ROBIN"}

	update := Params{"test1", "LEAST_CONNECTIONS"}

	fullName := "gcore_lbpool.acctest"

	ripTemplate := func(params *Params) string {
		return fmt.Sprintf(`
            resource "gcore_lbpool" "acctest" {
			  %s
              %s
			  name = "%s"
			  protocol = "HTTP"
			  lb_algorithm = "%s"
			  loadbalancer_id = "%s"
			  listener_id = "%s"
			}
		`, projectInfo(), regionInfo(), params.Name, params.LBAlgorithm, lbID, listener.ID)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccLBPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: ripTemplate(&create),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(fullName),
					resource.TestCheckResourceAttr(fullName, "name", create.Name),
					resource.TestCheckResourceAttr(fullName, "lb_algorithm", create.LBAlgorithm),
				),
			},
			{
				Config: ripTemplate(&update),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(fullName),
					resource.TestCheckResourceAttr(fullName, "name", update.Name),
					resource.TestCheckResourceAttr(fullName, "lb_algorithm", update.LBAlgorithm),
				),
			},
		},
	})
}

func testAccLBPoolDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	client, err := CreateTestClient(config.Provider, LBPoolsPoint, versionPointV1)
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "gcore_lbpool" {
			_, err := lbpools.Get(client, rs.Primary.ID).Extract()
			if err == nil {
				return fmt.Errorf("LBPool still exists")
			}
		}
	}

	return nil
}
