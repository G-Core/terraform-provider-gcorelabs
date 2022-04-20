provider gcore {
  permanent_api_token = "251$d3361.............1b35f26d8"
}

//
// example0: managing zone and records by TF using variables
//
variable "example_domain0" {
  type    = string
  default = "examplezone.com"
}

resource "gcore_dns_zone" "examplezone0" {
  name = var.example_domain0
}

resource "gcore_dns_zone_record" "example_rrset0" {
  zone   = gcore_dns_zone.examplezone0.name
  domain = "${gcore_dns_zone.examplezone0.name}"
  type   = "A"
  ttl    = 100

  resource_record {
    content = "127.0.0.100"
  }
  resource_record {
    content = "127.0.0.200"
    // enabled = false
  }
}

//
// example1: managing zone outside of TF 
//
resource "gcore_dns_zone_record" "subdomain_examplezone" {
  zone   = "examplezone.com"
  domain = "subdomain.examplezone.com"
  type   = "TXT"
  ttl    = 10

  filter {
    type   = "geodistance"
    limit  = 1
    strict = true
  }

  resource_record {
    content = "1234"
    enabled = true

    meta {
      latlong    = [52.367, 4.9041]
      asn        = [12345]
      ip         = ["1.1.1.1"]
      notes      = ["notes"]
      continents = ["asia"]
      countries  = ["russia"]
      default    = true
    }
  }
}

resource "gcore_dns_zone_record" "subdomain_examplezone_mx" {
  zone   = "examplezone.com"
  domain = "subdomain.examplezone.com"
  type   = "MX"
  ttl    = 10

  resource_record {
    content = "10 mail.my.com."
    enabled = true
  }
}

resource "gcore_dns_zone_record" "subdomain_examplezone_caa" {
  zone   = "examplezone.com"
  domain = "subdomain.examplezone.com"
  type   = "CAA"
  ttl    = 10

  resource_record {
    content = "0 issue \"company.org; account=12345\""
    enabled = true
  }
}
