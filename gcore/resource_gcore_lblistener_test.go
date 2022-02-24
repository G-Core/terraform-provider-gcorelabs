package gcore

import (
	"fmt"
	"testing"

	"github.com/G-Core/gcorelabscloud-go/gcore/loadbalancer/v1/listeners"
	"github.com/G-Core/gcorelabscloud-go/gcore/loadbalancer/v1/loadbalancers"
	"github.com/G-Core/gcorelabscloud-go/gcore/loadbalancer/v1/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccLBListener(t *testing.T) {
	cfg, err := createTestConfig()
	if err != nil {
		t.Fatal(err)
	}

	client, err := CreateTestClient(cfg.Provider, LoadBalancersPoint, versionPointV1)
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

	type Params struct {
		Name string
	}

	create := Params{"test"}

	update := Params{"test1"}

	fullName := "gcore_lblistener.acctest"

	ripTemplate := func(params *Params) string {
		return fmt.Sprintf(`
            resource "gcore_lblistener" "acctest" {
			  %s
              %s
			  name = "%s"
			  protocol = "TCP"
			  protocol_port = 36621
			  loadbalancer_id = "%s"
			}
		`, projectInfo(), regionInfo(), params.Name, lbID)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccLBListenerDestroy,
		Steps: []resource.TestStep{
			{
				Config: ripTemplate(&create),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(fullName),
					resource.TestCheckResourceAttr(fullName, "name", create.Name),
				),
			},
			{
				Config: ripTemplate(&update),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(fullName),
					resource.TestCheckResourceAttr(fullName, "name", update.Name),
				),
			},
		},
	})
}

func testAccLBListenerDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	client, err := CreateTestClient(config.Provider, LBListenersPoint, versionPointV1)
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "gcore_lblistener" {
			_, err := listeners.Get(client, rs.Primary.ID).Extract()
			if err == nil {
				return fmt.Errorf("LBListener still exists")
			}
		}
	}

	return nil
}
