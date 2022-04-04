//go:build cloud
// +build cloud

package gcore

import (
	"fmt"
	"testing"

	"github.com/G-Core/gcorelabscloud-go/gcore/secret/v1/secrets"
	secretsV2 "github.com/G-Core/gcorelabscloud-go/gcore/secret/v2/secrets"
	"github.com/G-Core/gcorelabscloud-go/gcore/task/v1/tasks"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccSecretDataSource(t *testing.T) {
	cfg, err := createTestConfig()
	if err != nil {
		t.Fatal(err)
	}

	client, err := CreateTestClient(cfg.Provider, secretPoint, versionPointV1)
	if err != nil {
		t.Fatal(err)
	}
	clientV2, err := CreateTestClient(cfg.Provider, secretPoint, versionPointV2)
	if err != nil {
		t.Fatal(err)
	}

	opts := secretsV2.CreateOpts{
		Name: secretName,
		Payload: secretsV2.PayloadOpts{
			CertificateChain: certificateChain,
			Certificate:      certificate,
			PrivateKey:       privateKey,
		},
	}
	results, err := secretsV2.Create(clientV2, opts).Extract()
	if err != nil {
		t.Fatal(err)
	}

	taskID := results.Tasks[0]
	secretID, err := tasks.WaitTaskAndReturnResult(client, taskID, true, SecretCreatingTimeout, func(task tasks.TaskID) (interface{}, error) {
		taskInfo, err := tasks.Get(client, string(task)).Extract()
		if err != nil {
			return nil, fmt.Errorf("cannot get task with ID: %s. Error: %w", task, err)
		}
		Secret, err := secrets.ExtractSecretIDFromTask(taskInfo)
		if err != nil {
			return nil, fmt.Errorf("cannot retrieve Secret ID from task info: %w", err)
		}
		return Secret, nil
	},
	)

	if err != nil {
		t.Fatal(err)
	}
	defer secrets.Delete(client, secretID.(string))

	fullName := "data.gcore_secret.acctest"
	kpTemplate := fmt.Sprintf(`
	data "gcore_secret" "acctest" {
	  %s
      %s
      name = "%s"
	}
	`, projectInfo(), regionInfo(), secretName)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: kpTemplate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(fullName),
					resource.TestCheckResourceAttr(fullName, "name", secretName),
				),
			},
		},
	})
}
