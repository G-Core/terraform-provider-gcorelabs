package gcore

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccK8sPoolDataSource(t *testing.T) {
	fullName := "data.gcore_k8s_pool.acctest"

	ipTemplate := fmt.Sprintf(`
			data "gcore_k8s_pool" "acctest" {
			  %s
              %s
              cluster_id = "%s"
			  pool_id = "%s"
			}
		`, projectInfo(), regionInfo(), GCORE_CLUSTER_ID, GCORE_CLUSTER_POOL_ID)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckK8sPoolDataSource(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: ipTemplate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(fullName),
					resource.TestCheckResourceAttr(fullName, "cluster_id", GCORE_CLUSTER_ID),
					resource.TestCheckResourceAttr(fullName, "pool_id", GCORE_CLUSTER_POOL_ID),
				),
			},
		},
	})
}
