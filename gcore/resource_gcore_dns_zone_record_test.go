package gcore

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDnsZoneRecord(t *testing.T) {

	random := time.Now().Nanosecond()
	subDomain := "terraformtest"
	domain := fmt.Sprintf("key%d", random)
	name := fmt.Sprintf("%s_%s", subDomain, domain)
	zone := domain + ".com"
	fullDomain := subDomain + "." + zone

	resourceName := fmt.Sprintf("%s.%s", DNSZoneRecordResource, name)

	templateCreate := func() string {
		return fmt.Sprintf(`
resource "%s" "%s" {
  zone = "%s"
  domain = "%s"
  type = "TXT"
  ttl = 10

  resource_records {
    content  = "1234"
    
    meta {
      latlong = [52.367,4.9041]
	  asn = [12345]
	  ip = ["1.1.1.1"]
	  notes = ["notes"]
	  continents = ["asia"]
	  countries = ["russia"]
	  default = true
  	}
  }
}
		`, DNSZoneRecordResource, name, zone, fullDomain)
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckVars(t, GCORE_USERNAME_VAR, GCORE_PASSWORD_VAR, GCORE_DNS_URL_VAR)
		},
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: templateCreate(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, DNSZoneRecordSchemaDomain, fullDomain),
					resource.TestCheckResourceAttr(resourceName, DNSZoneRecordSchemaType, "TXT"),
					resource.TestCheckResourceAttr(resourceName, DNSZoneRecordSchemaTTL, "10"),
					resource.TestCheckResourceAttr(resourceName,
						fmt.Sprintf("%s.0.%s", DNSZoneRecordSchemaResourceRecords, DNSZoneRecordSchemaContent),
						"1234"),
					resource.TestCheckResourceAttr(resourceName,
						fmt.Sprintf("%s.0.%s.0.%s.0",
							DNSZoneRecordSchemaResourceRecords, DNSZoneRecordSchemaMeta, DNSZoneRecordSchemaMetaLatLong),
						"52.367"),
					resource.TestCheckResourceAttr(resourceName,
						fmt.Sprintf("%s.0.%s.0.%s.1",
							DNSZoneRecordSchemaResourceRecords, DNSZoneRecordSchemaMeta, DNSZoneRecordSchemaMetaLatLong),
						"4.9041"),
					resource.TestCheckResourceAttr(resourceName,
						fmt.Sprintf("%s.0.%s.0.%s.0",
							DNSZoneRecordSchemaResourceRecords, DNSZoneRecordSchemaMeta, DNSZoneRecordSchemaMetaAsn),
						"12345"),
					resource.TestCheckResourceAttr(resourceName,
						fmt.Sprintf("%s.0.%s.0.%s.0",
							DNSZoneRecordSchemaResourceRecords, DNSZoneRecordSchemaMeta, DNSZoneRecordSchemaMetaIP),
						"1.1.1.1"),
					resource.TestCheckResourceAttr(resourceName,
						fmt.Sprintf("%s.0.%s.0.%s.0",
							DNSZoneRecordSchemaResourceRecords, DNSZoneRecordSchemaMeta, DNSZoneRecordSchemaMetaNotes),
						"notes"),
					resource.TestCheckResourceAttr(resourceName,
						fmt.Sprintf("%s.0.%s.0.%s.0",
							DNSZoneRecordSchemaResourceRecords, DNSZoneRecordSchemaMeta, DNSZoneRecordSchemaMetaContinents),
						"asia"),
					resource.TestCheckResourceAttr(resourceName,
						fmt.Sprintf("%s.0.%s.0.%s.0",
							DNSZoneRecordSchemaResourceRecords, DNSZoneRecordSchemaMeta, DNSZoneRecordSchemaMetaCountries),
						"russia"),
					resource.TestCheckResourceAttr(resourceName,
						fmt.Sprintf("%s.0.%s.0.%s",
							DNSZoneRecordSchemaResourceRecords, DNSZoneRecordSchemaMeta, DNSZoneRecordSchemaMetaDefault),
						"true"),
				),
			},
		},
	})
}
