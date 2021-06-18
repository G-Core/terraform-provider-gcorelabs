package gcore

import (
	"fmt"
	"testing"

	"github.com/G-Core/gcorelabscloud-go/gcore/k8s/v1/clusters"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccK8s(t *testing.T) {
	fullName := "gcore_k8s.acctest"

	ipTemplate := fmt.Sprintf(`
			resource "gcore_k8s" "acctest" {
			  %s
              %s
              name = "tf-k8s"
			  fixed_network = "%s"
			  fixed_subnet = "%s"
			  pool {
				name = "tf-pool1"
				flavor_id = "g1-standard-1-2"
				min_node_count = 1
				max_node_count = 1
				node_count = 1
				docker_volume_size = 2
			  }

			}
		`, projectInfo(), regionInfo(), GCORE_NETWORK_ID, GCORE_SUBNET_ID)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckK8s(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccK8sDestroy,
		Steps: []resource.TestStep{
			{
				Config: ipTemplate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(fullName),
					resource.TestCheckResourceAttr(fullName, "name", "tf-k8s"),
				),
			},
		},
	})
}

func testAccK8sDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	client, err := CreateTestClient(config.Provider, K8sPoint, versionPointV1)
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gcore_k8s" {
			continue
		}

		_, err := clusters.Get(client, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("k8s cluster still exists")
		}
	}

	return nil
}
