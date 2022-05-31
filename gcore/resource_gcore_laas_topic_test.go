//go:build cloud
// +build cloud

package gcore

import (
	"fmt"
	"testing"

	"github.com/G-Core/gcorelabscloud-go/gcore/laas/v1/laas"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const (
	topicName = "test-topic"
)

func TestAccLaaSTopic(t *testing.T) {
	fullName := "gcore_laas_topic.acctest"
	kpTemplate := fmt.Sprintf(`
	resource "gcore_laas_topic" "acctest" {
	  %s
      %s
      name = "%s"
	}
	`, projectInfo(), regionInfo(), topicName)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccLaaSTopicDestroy,
		Steps: []resource.TestStep{
			{
				Config: kpTemplate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(fullName),
					resource.TestCheckResourceAttr(fullName, "name", topicName),
				),
			},
		},
	})
}

func testAccLaaSTopicDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	client, err := CreateTestClient(config.Provider, laasPoint, versionPointV1)
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gcore_laas_topic" {
			continue
		}

		topics, err := laas.ListTopicAll(client)
		if err != nil {
			return fmt.Errorf("cant get topic list: %s", err.Error())
		}
		for _, t := range topics {
			if t.Name == topicName {
				return fmt.Errorf("secret still exists")
			}
		}
		return nil
	}

	return nil
}
