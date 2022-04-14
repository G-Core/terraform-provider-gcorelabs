provider gcore {
  user_name = "test"
  password = "test"
  permanent_api_token="123$321" // https://support.gcorelabs.com/hc/en-us/articles/360018625617-API-tokens
  ignore_creds_auth_error=true // if you want to manage storage or dns resources only and provide permanent_api_token without user_name & password
  gcore_platform_api = "https://api.gcorelabs.com"
  gcore_dns_api = "https://api.gcorelabs.com/dns"
}

resource "gcore_dns_zone" "example_zone" {
  name = "example_zone.com"
}
