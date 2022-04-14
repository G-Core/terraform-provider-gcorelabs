provider gcore {
  user_name = "test"
  password = "test"
  permanent_api_token="123$321" // https://support.gcorelabs.com/hc/en-us/articles/360018625617-API-tokens
  ignore_creds_auth_error=true // if you want to manage storage resource only and provide permanent_api_token without user_name & password
  gcore_platform_api = "https://api.gcorelabs.com"
  gcore_storage_api = "https://storage.gcorelabs.com/api"
}

data "gcore_storage_sftp" "example_sftp" {
  name = "example"
}
