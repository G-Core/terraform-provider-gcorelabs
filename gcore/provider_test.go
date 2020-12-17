package gcore

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"os"
	"strconv"
	"strings"
	"testing"

	gcorecloud "github.com/G-Core/gcorelabscloud-go"
	gc "github.com/G-Core/gcorelabscloud-go/gcore"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	GCORE_USERNAME    = os.Getenv("GCORE_USERNAME")
	GCORE_PASSWORD    = os.Getenv("GCORE_PASSWORD")
	GCORE_EXT_NET     = os.Getenv("GCORE_EXT_NET")
	GCORE_PRIV_SUBNET = os.Getenv("GCORE_PRIV_SUBNET")
)

var testAccProvider *schema.Provider
var testAccProviders map[string]func() (*schema.Provider, error)

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]func() (*schema.Provider, error){
		"gcore": func() (*schema.Provider, error) {
			return testAccProvider, nil
		},
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func testAccPreCheck(t *testing.T) {
	Vars := map[string]interface{}{
		"GCORE_USERNAME": GCORE_USERNAME,
		"GCORE_PASSWORD": GCORE_PASSWORD,
	}
	for k, v := range Vars {
		if v == "" {
			t.Fatalf("'%s' must be set for acceptance test", k)
		}
	}
	checkNameAndID("PROJECT", t)
	checkNameAndID("REGION", t)
}

func testAccPreCheckRouter(t *testing.T) {
	Vars := map[string]interface{}{
		"GCORE_EXT_NET":     GCORE_EXT_NET,
		"GCORE_PRIV_SUBNET": GCORE_PRIV_SUBNET,
	}
	for k, v := range Vars {
		if v == "" {
			t.Fatalf("'%s' must be set for acceptance test", k)
		}
	}
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

func CreateTestClient(provider *gcorecloud.ProviderClient, endpoint string) (*gcorecloud.ServiceClient, error) {
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

	client, err := gc.ClientServiceFromProvider(provider, gcorecloud.EndpointOpts{
		Name:    endpoint,
		Region:  regionID,
		Project: projectID,
		Version: "v1",
	})

	if err != nil {
		return nil, err
	}
	return client, nil
}

func testAccCheckResourceExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// retrieve the resource by name from state
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Widget ID is not set")
		}
		return nil
	}
}
