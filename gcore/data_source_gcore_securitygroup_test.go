package gcore

import (
	"fmt"
	"testing"

	"github.com/G-Core/gcorelabscloud-go/gcore/securitygroup/v1/securitygroups"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const (
	securityGroupTestName = "test-sg"
)

func TestAccSecurityGroupDataSource(t *testing.T) {
	cfg, err := createTestConfig()
	if err != nil {
		t.Fatal(err)
	}

	client, err := CreateTestClient(cfg.Provider, securityGroupPoint, versionPointV1)
	if err != nil {
		t.Fatal(err)
	}

	opts := securitygroups.CreateOpts{
		SecurityGroup: securitygroups.CreateSecurityGroupOpts{
			Name:               securityGroupTestName,
			SecurityGroupRules: []securitygroups.CreateSecurityGroupRuleOpts{},
		},
	}

	sg, err := securitygroups.Create(client, opts).Extract()
	if err != nil {
		t.Fatal(err)
	}

	defer securitygroups.Delete(client, sg.ID)

	fullName := "data.gcore_securitygroup.acctest"
	tpl := func(name string) string {
		return fmt.Sprintf(`
			data "gcore_securitygroup" "acctest" {
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
				Config: tpl(sg.Name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(fullName),
					resource.TestCheckResourceAttr(fullName, "name", sg.Name),
					resource.TestCheckResourceAttr(fullName, "id", sg.ID),
				),
			},
		},
	})
}
