package gcore

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDnsZoneRecord(t *testing.T) {

	random := time.Now().Nanosecond()
	domain := "terraformtest"
	subDomain := fmt.Sprintf("key%d", random)
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

  filter {
    type = "geodistance"
    limit = 1
    strict = true
  }

  resource_record {
    content  = "1234"
    enabled = true
    
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
	templateUpdate := func() string {
		return fmt.Sprintf(`
resource "%s" "%s" {
  zone = "%s"
  domain = "%s"
  type = "TXT"
  ttl = 20

  resource_record {
    content  = "12345"
    
    meta {
      latlong = [52.367,4.9041]
	  ip = ["1.1.2.2"]
	  notes = ["notes"]
	  continents = ["america"]
	  countries = ["usa"]
	  default = false
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
						fmt.Sprintf("%s.0.%s", DNSZoneRecordSchemaFilter, DNSZoneRecordSchemaFilterType),
						"geodistance"),
					resource.TestCheckResourceAttr(resourceName,
						fmt.Sprintf("%s.0.%s", DNSZoneRecordSchemaFilter, DNSZoneRecordSchemaFilterLimit),
						"1"),
					resource.TestCheckResourceAttr(resourceName,
						fmt.Sprintf("%s.0.%s", DNSZoneRecordSchemaFilter, DNSZoneRecordSchemaFilterStrict),
						"true"),
					resource.TestCheckResourceAttr(resourceName,
						fmt.Sprintf("%s.0.%s", DNSZoneRecordSchemaResourceRecord, DNSZoneRecordSchemaContent),
						"1234"),
					resource.TestCheckResourceAttr(resourceName,
						fmt.Sprintf("%s.0.%s", DNSZoneRecordSchemaResourceRecord, DNSZoneRecordSchemaEnabled),
						"true"),
					resource.TestCheckResourceAttr(resourceName,
						fmt.Sprintf("%s.0.%s.0.%s.0",
							DNSZoneRecordSchemaResourceRecord, DNSZoneRecordSchemaMeta, DNSZoneRecordSchemaMetaLatLong),
						"52.367"),
					resource.TestCheckResourceAttr(resourceName,
						fmt.Sprintf("%s.0.%s.0.%s.1",
							DNSZoneRecordSchemaResourceRecord, DNSZoneRecordSchemaMeta, DNSZoneRecordSchemaMetaLatLong),
						"4.9041"),
					resource.TestCheckResourceAttr(resourceName,
						fmt.Sprintf("%s.0.%s.0.%s.0",
							DNSZoneRecordSchemaResourceRecord, DNSZoneRecordSchemaMeta, DNSZoneRecordSchemaMetaAsn),
						"12345"),
					resource.TestCheckResourceAttr(resourceName,
						fmt.Sprintf("%s.0.%s.0.%s.0",
							DNSZoneRecordSchemaResourceRecord, DNSZoneRecordSchemaMeta, DNSZoneRecordSchemaMetaIP),
						"1.1.1.1"),
					resource.TestCheckResourceAttr(resourceName,
						fmt.Sprintf("%s.0.%s.0.%s.0",
							DNSZoneRecordSchemaResourceRecord, DNSZoneRecordSchemaMeta, DNSZoneRecordSchemaMetaNotes),
						"notes"),
					resource.TestCheckResourceAttr(resourceName,
						fmt.Sprintf("%s.0.%s.0.%s.0",
							DNSZoneRecordSchemaResourceRecord, DNSZoneRecordSchemaMeta, DNSZoneRecordSchemaMetaContinents),
						"asia"),
					resource.TestCheckResourceAttr(resourceName,
						fmt.Sprintf("%s.0.%s.0.%s.0",
							DNSZoneRecordSchemaResourceRecord, DNSZoneRecordSchemaMeta, DNSZoneRecordSchemaMetaCountries),
						"russia"),
					resource.TestCheckResourceAttr(resourceName,
						fmt.Sprintf("%s.0.%s.0.%s",
							DNSZoneRecordSchemaResourceRecord, DNSZoneRecordSchemaMeta, DNSZoneRecordSchemaMetaDefault),
						"true"),
				),
			},
			{
				Config: templateUpdate(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, DNSZoneRecordSchemaDomain, fullDomain),
					resource.TestCheckResourceAttr(resourceName, DNSZoneRecordSchemaType, "TXT"),
					resource.TestCheckResourceAttr(resourceName, DNSZoneRecordSchemaTTL, "20"),
					resource.TestCheckResourceAttr(resourceName,
						fmt.Sprintf("%s.0.%s", DNSZoneRecordSchemaResourceRecord, DNSZoneRecordSchemaContent),
						"12345"),
					resource.TestCheckResourceAttr(resourceName,
						fmt.Sprintf("%s.0.%s.0.%s.0",
							DNSZoneRecordSchemaResourceRecord, DNSZoneRecordSchemaMeta, DNSZoneRecordSchemaMetaLatLong),
						"52.367"),
					resource.TestCheckResourceAttr(resourceName,
						fmt.Sprintf("%s.0.%s.0.%s.1",
							DNSZoneRecordSchemaResourceRecord, DNSZoneRecordSchemaMeta, DNSZoneRecordSchemaMetaLatLong),
						"4.9041"),
					resource.TestCheckNoResourceAttr(resourceName, fmt.Sprintf("%s.0.%s.0.%s.0",
						DNSZoneRecordSchemaResourceRecord, DNSZoneRecordSchemaMeta, DNSZoneRecordSchemaMetaAsn)),
					resource.TestCheckResourceAttr(resourceName,
						fmt.Sprintf("%s.0.%s.0.%s.0",
							DNSZoneRecordSchemaResourceRecord, DNSZoneRecordSchemaMeta, DNSZoneRecordSchemaMetaIP),
						"1.1.2.2"),
					resource.TestCheckResourceAttr(resourceName,
						fmt.Sprintf("%s.0.%s.0.%s.0",
							DNSZoneRecordSchemaResourceRecord, DNSZoneRecordSchemaMeta, DNSZoneRecordSchemaMetaNotes),
						"notes"),
					resource.TestCheckResourceAttr(resourceName,
						fmt.Sprintf("%s.0.%s.0.%s.0",
							DNSZoneRecordSchemaResourceRecord, DNSZoneRecordSchemaMeta, DNSZoneRecordSchemaMetaContinents),
						"america"),
					resource.TestCheckResourceAttr(resourceName,
						fmt.Sprintf("%s.0.%s.0.%s.0",
							DNSZoneRecordSchemaResourceRecord, DNSZoneRecordSchemaMeta, DNSZoneRecordSchemaMetaCountries),
						"usa"),
					resource.TestCheckResourceAttr(resourceName,
						fmt.Sprintf("%s.0.%s.0.%s",
							DNSZoneRecordSchemaResourceRecord, DNSZoneRecordSchemaMeta, DNSZoneRecordSchemaMetaDefault),
						"false"),
				),
			},
		},
	})
}
