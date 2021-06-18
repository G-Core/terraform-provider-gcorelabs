package gcore

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/G-Core/gcorelabscloud-go/gcore/k8s/v1/pools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccK8sPool(t *testing.T) {
	fullName := "gcore_k8s_pool.acctest"

	type Params struct {
		Name             string
		Flavor           string
		MinNodeCount     int
		MaxNodeCount     int
		NodeCount        int
		DockerVolumeSize int
	}

	create := Params{
		Name:             "tf-pool1",
		Flavor:           "g1-standard-1-2",
		MinNodeCount:     1,
		MaxNodeCount:     1,
		NodeCount:        1,
		DockerVolumeSize: 2,
	}

	update := Params{
		Name:             "tf-pool2",
		Flavor:           "g1-standard-1-2",
		MinNodeCount:     1,
		MaxNodeCount:     2,
		NodeCount:        1,
		DockerVolumeSize: 2,
	}

	ipTemplate := func(p *Params) string {
		return fmt.Sprintf(`
			resource "gcore_k8s_pool" "acctest" {
			  %s
              %s
              cluster_id = "%s"
			  name = "%s"
			  flavor_id = "%s"
			  min_node_count = %d
			  max_node_count = %d
			  node_count = %d
			  docker_volume_size = %d
			}
		`, projectInfo(), regionInfo(), GCORE_CLUSTER_ID,
			p.Name, p.Flavor, p.MinNodeCount, p.MaxNodeCount,
			p.NodeCount, p.DockerVolumeSize)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckK8sPool(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccK8sPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: ipTemplate(&create),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(fullName),
					resource.TestCheckResourceAttr(fullName, "name", create.Name),
					resource.TestCheckResourceAttr(fullName, "flavor_id", create.Flavor),
					resource.TestCheckResourceAttr(fullName, "docker_volume_size", strconv.Itoa(create.DockerVolumeSize)),
					resource.TestCheckResourceAttr(fullName, "min_node_count", strconv.Itoa(create.MinNodeCount)),
					resource.TestCheckResourceAttr(fullName, "max_node_count", strconv.Itoa(create.MaxNodeCount)),
					resource.TestCheckResourceAttr(fullName, "node_count", strconv.Itoa(create.NodeCount)),
				),
			},
			{
				Config: ipTemplate(&update),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(fullName),
					resource.TestCheckResourceAttr(fullName, "name", update.Name),
					resource.TestCheckResourceAttr(fullName, "flavor_id", update.Flavor),
					resource.TestCheckResourceAttr(fullName, "docker_volume_size", strconv.Itoa(update.DockerVolumeSize)),
					resource.TestCheckResourceAttr(fullName, "min_node_count", strconv.Itoa(update.MinNodeCount)),
					resource.TestCheckResourceAttr(fullName, "max_node_count", strconv.Itoa(update.MaxNodeCount)),
					resource.TestCheckResourceAttr(fullName, "node_count", strconv.Itoa(update.NodeCount)),
				),
			},
		},
	})
}

func testAccK8sPoolDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	client, err := CreateTestClient(config.Provider, K8sPoint, versionPointV1)
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gcore_k8s_pool" {
			continue
		}

		_, err := pools.Get(client, GCORE_CLUSTER_ID, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("k8s pool still exists")
		}
	}

	return nil
}
