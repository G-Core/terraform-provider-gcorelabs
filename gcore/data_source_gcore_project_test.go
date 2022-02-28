//go:build cloud
// +build cloud

package gcore

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/G-Core/gcorelabscloud-go/gcore/project/v1/projects"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccProjectDataSource(t *testing.T) {
	cfg, err := createTestConfig()
	if err != nil {
		t.Fatal(err)
	}

	client, err := CreateTestClient(cfg.Provider, projectPoint, versionPointV1)
	if err != nil {
		t.Fatal(err)
	}

	prjs, err := projects.ListAll(client)
	if err != nil {
		t.Fatal(err)
	}

	if len(prjs) == 0 {
		t.Fatal("projects not found")
	}

	project := prjs[0]

	fullName := "data.gcore_project.acctest"
	tpl := func(name string) string {
		return fmt.Sprintf(`
			data "gcore_project" "acctest" {
              name = "%s"
			}
		`, name)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: tpl(project.Name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(fullName),
					resource.TestCheckResourceAttr(fullName, "name", project.Name),
					resource.TestCheckResourceAttr(fullName, "id", strconv.Itoa(project.ID)),
				),
			},
		},
	})
}
