package gcore

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccK8sDataSource(t *testing.T) {
	fullName := "data.gcore_k8s.acctest"

	ipTemplate := fmt.Sprintf(`
			data "gcore_k8s" "acctest" {
			  %s
              %s
              cluster_id = "%s"
			}
		`, projectInfo(), regionInfo(), GCORE_CLUSTER_ID)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckK8sDataSource(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: ipTemplate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(fullName),
					resource.TestCheckResourceAttr(fullName, "cluster_id", GCORE_CLUSTER_ID),
				),
			},
		},
	})
}
