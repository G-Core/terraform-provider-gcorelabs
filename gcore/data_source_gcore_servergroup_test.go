package gcore

import (
	"fmt"
	"testing"

	"github.com/G-Core/gcorelabscloud-go/gcore/servergroup/v1/servergroups"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccServerGroupDataSource(t *testing.T) {
	cfg, err := createTestConfig()
	if err != nil {
		t.Fatal(err)
	}

	client, err := CreateTestClient(cfg.Provider, serverGroupsPoint, versionPointV1)
	if err != nil {
		t.Fatal(err)
	}

	opts := servergroups.CreateOpts{Name: "name", Policy: servergroups.AntiAffinityPolicy}
	serverGroup, err := servergroups.Create(client, opts).Extract()
	if err != nil {
		t.Fatal(err)
	}

	fullName := "data.gcore_servergroup.acctest"
	tpl := func(name string) string {
		return fmt.Sprintf(`
			data "gcore_servergroup" "acctest" {
			  %s
              %s
              name = "%s"
			}
		`, projectInfo(), regionInfo(), name)
	}

	defer servergroups.Delete(client, serverGroup.ServerGroupID)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: tpl(opts.Name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(fullName),
					resource.TestCheckResourceAttr(fullName, "name", serverGroup.Name),
					resource.TestCheckResourceAttr(fullName, "id", serverGroup.ServerGroupID),
					resource.TestCheckResourceAttr(fullName, "policy", serverGroup.Policy.String()),
				),
			},
		},
	})
}
