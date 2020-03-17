package main

import (
	"fmt"
	"os"
	"strings"
	"testing"

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
	if os.Getenv("OS_PROVIDER_JWT") == "" {
		t.Fatal("OS_PROVIDER_JWT must be set for acceptance tests")
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

func providerData() string {
	return fmt.Sprintf(`
	provider "gcore" {
		jwt = "%s"
	}`, os.Getenv("TEST_PROVIDER_JWT"))
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
