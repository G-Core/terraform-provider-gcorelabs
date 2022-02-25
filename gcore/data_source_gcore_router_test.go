//go:build cloud
// +build cloud

package gcore

import (
	"fmt"
	"testing"

	"github.com/G-Core/gcorelabscloud-go/gcore/network/v1/networks"
	"github.com/G-Core/gcorelabscloud-go/gcore/router/v1/routers"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccRouterDataSource(t *testing.T) {
	cfg, err := createTestConfig()
	if err != nil {
		t.Fatal(err)
	}

	clientNet, err := CreateTestClient(cfg.Provider, networksPoint, versionPointV1)
	if err != nil {
		t.Fatal(err)
	}

	clientRouter, err := CreateTestClient(cfg.Provider, RouterPoint, versionPointV1)
	if err != nil {
		t.Fatal(err)
	}

	opts := networks.CreateOpts{
		Name:         networkTestName,
		CreateRouter: true,
	}

	networkID, err := createTestNetwork(clientNet, opts)
	if err != nil {
		t.Fatal(err)
	}
	defer networks.Delete(clientNet, networkID)

	rs, err := routers.ListAll(clientRouter, routers.ListOpts{})
	if err != nil {
		t.Fatal(err)
	}
	router := rs[0]

	fullName := "data.gcore_router.acctest"
	tpl := func(name string) string {
		return fmt.Sprintf(`
			data "gcore_router" "acctest" {
			  %s
              %s
              name = "%s"
			}
		`, projectInfo(), regionInfo(), name)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: tpl(router.Name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(fullName),
					resource.TestCheckResourceAttr(fullName, "name", router.Name),
					resource.TestCheckResourceAttr(fullName, "id", router.ID),
				),
			},
		},
	})
}
