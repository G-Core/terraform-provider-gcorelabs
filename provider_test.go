package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	"github.com/G-Core/gcorelabscloud-go/gcore"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"gcore": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func testAccPreCheck(t *testing.T) {
	if os.Getenv("GCORE_PROVIDER_USERNAME") == "" {
		t.Fatal("GCORE_PROVIDER_USERNAME must be set for acceptance tests")
	}
	if os.Getenv("GCORE_PROVIDER_PASSWORD") == "" {
		t.Fatal("GCORE_PROVIDER_PASSWORD must be set for acceptance tests")
	}
	checkNameAndID("PROJECT", t)
	checkNameAndID("REGION", t)
}

func checkNameAndID(resourceType string, t *testing.T) {
	// resourceType is a word in capital letters
	keyID := fmt.Sprintf("TEST_%s_ID", resourceType)
	keyNane := fmt.Sprintf("TEST_%s_NAME", resourceType)
	_, haveID := os.LookupEnv(keyID)
	_, haveName := os.LookupEnv(keyNane)
	if !haveID && !haveName {
		t.Fatalf("%s or %s must be set for acceptance tests", keyID, keyNane)
	}
	if haveID && haveName {
		t.Fatalf("Use only one from environment variables: %s or %s", keyID, keyNane)
	}
}

func regionInfo() string {
	return objectInfo("REGION")
}

func projectInfo() string {
	return objectInfo("PROJECT")
}

func objectInfo(resourceType string) string {
	// resourceType is a word in capital letters
	keyID := fmt.Sprintf("TEST_%s_ID", resourceType)
	keyNane := fmt.Sprintf("TEST_%s_NAME", resourceType)
	if regionID, exists := os.LookupEnv(keyID); exists {
		return fmt.Sprintf(`%s_id = %s`, strings.ToLower(resourceType), regionID)
	}
	return fmt.Sprintf(`%s_name = "%s"`, strings.ToLower(resourceType), os.Getenv(keyNane))
}

func CreateTestClient(provider *gcorecloud.ProviderClient) (*gcorecloud.ServiceClient, error) {
	projectID := 0
	err := fmt.Errorf("")
	if strProjectID, exists := os.LookupEnv("TEST_PROJECT_ID"); exists {
		projectID, err = strconv.Atoi(strProjectID)
		if err != nil {
			return nil, err
		}
	} else {
		projectID, err = GetProject(provider, 0, os.Getenv("TEST_PROJECT_NAME"))
		if err != nil {
			return nil, err
		}
	}
	regionID := 0
	if strRegionID, exists := os.LookupEnv("TEST_REGION_ID"); exists {
		regionID, err = strconv.Atoi(strRegionID)
		if err != nil {
			return nil, err
		}
	} else {
		regionID, err = GetProject(provider, 0, os.Getenv("TEST_REGION_NAME"))
		if err != nil {
			return nil, err
		}
	}

	client, err := gcore.ClientServiceFromProvider(provider, gcorecloud.EndpointOpts{
		Name:    "volumes",
		Region:  regionID,
		Project: projectID,
		Version: "v1",
	})

	if err != nil {
		return nil, err
	}
	return client, nil
}
