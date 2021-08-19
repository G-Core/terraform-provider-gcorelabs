provider gcore {
  user_name = "test"
  password = "test"
  permanent_api_token="123$321" // https://support.gcorelabs.com/hc/en-us/articles/360018625617-API-tokens
  ignore_creds_auth_error=true // if you want to manage storage or dns resources only and provide permanent_api_token without user_name & password
  gcore_platform = "https://api.gcdn.co"
  gcore_dns_api = "https://dnsapi.gcorelabs.com"
}

resource "gcore_dns_zone_record" "subdomain_examplezone" {
  zone = "examplezone.com"
  domain = "subdomain.examplezone.com"
  type = "TXT"
  ttl = 10

  resource_record {
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
